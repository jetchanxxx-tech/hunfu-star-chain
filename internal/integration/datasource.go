package integration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/huifu/star-chain/internal/model"
	"github.com/huifu/star-chain/internal/store"
)

type DataSourceType string

const (
	HIS  DataSourceType = "his"
	LIS  DataSourceType = "lis"
	PACS DataSourceType = "pacs"
	EMR  DataSourceType = "emr"
)

type DataSourceConfig struct {
	Code       string         `json:"code"`
	Name       string         `json:"name"`
	SourceType DataSourceType `json:"source_type"`
	ConnMode   string         `json:"conn_mode"` // direct/ssl/vpn
	Host       string         `json:"host"`
	Port       int            `json:"port"`
	DBName     string         `json:"db_name"`
	Username   string         `json:"username"`
	Password   string         `json:"password"`
	QuerySQL   string         `json:"query_sql"` // Custom SQL / stored procedure
}

type FieldMapping struct {
	SourceField string `json:"source_field"` // Hospital field name
	TargetField string `json:"target_field"` // Platform field name
	Transform   string `json:"transform"`    // Transformation rule: direct/lab_map/date_parse/json_extract
	TransformArg string `json:"transform_arg"`
}

type MappingConfig struct {
	SourceCode DataSourceType `json:"source_code"`
	Mappings   []FieldMapping `json:"mappings"`
}

// SyncManager manages data sync from hospital systems
type SyncManager struct {
	db          *sql.DB
	reportStore *store.ReportStore
	mappings    map[string]*MappingConfig
	mu          sync.RWMutex
}

func NewSyncManager(db *sql.DB) *SyncManager {
	return &SyncManager{
		db:          db,
		reportStore: store.NewReportStore(db),
		mappings:    make(map[string]*MappingConfig),
	}
}

func (m *SyncManager) RegisterMapping(code string, cfg *MappingConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.mappings[code] = cfg
}

func (m *SyncManager) GetMapping(code string) *MappingConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.mappings[code]
}

// SyncFromSource pulls data from a hospital data source
func (m *SyncManager) SyncFromSource(cfg DataSourceConfig) (*model.DataSyncLog, error) {
	syncLog := &model.DataSyncLog{
		Source:   cfg.Code,
		SyncType: "incremental",
		StartedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	// Open connection to hospital data source
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?readTimeout=30s&timeout=10s&charset=utf8mb4",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	hospitalDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to hospital db: %w", err)
	}
	defer hospitalDB.Close()
	hospitalDB.SetConnMaxLifetime(30 * time.Second)
	hospitalDB.SetMaxOpenConns(2)

	rows, err := hospitalDB.Query(cfg.QuerySQL)
	if err != nil {
		return nil, fmt.Errorf("query hospital: %w", err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	mapping := m.GetMapping(cfg.Code)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			syncLog.FailCount++
			m.writeToDeadLetter(cfg.Code, fmt.Sprintf("%v", values), err.Error())
			continue
		}

		record := make(map[string]interface{})
		for i, col := range columns {
			record[col] = values[i]
		}

		// Apply field mapping
		report, err := m.mapToReport(record, mapping)
		if err != nil {
			syncLog.FailCount++
			rawJSON, _ := json.Marshal(record)
			m.writeToDeadLetter(cfg.Code, string(rawJSON), err.Error())
			continue
		}

		// Dedup: check report_no exists
		if report.ReportNo.Valid {
			existing, _ := m.reportStore.FindByReportNo(report.ReportNo.String)
			if existing != nil {
				// Update if newer
				continue
			}
		}

		if err := m.reportStore.Insert(report); err != nil {
			syncLog.FailCount++
			rawJSON, _ := json.Marshal(record)
			m.writeToDeadLetter(cfg.Code, string(rawJSON), err.Error())
			continue
		}
		syncLog.TotalCount++
		syncLog.SuccessCount++
	}

	syncLog.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}
	return syncLog, nil
}

func (m *SyncManager) mapToReport(record map[string]interface{}, mapping *MappingConfig) (*model.HealthReport, error) {
	report := &model.HealthReport{Source: "manual"}

	if mapping == nil {
		return nil, fmt.Errorf("no mapping config found")
	}

	fieldMap := make(map[string]interface{})
	for k, v := range record {
		fieldMap[k] = v
	}

	for _, fm := range mapping.Mappings {
		val, ok := fieldMap[fm.SourceField]
		if !ok {
			continue
		}
		valStr := fmt.Sprintf("%v", val)

		transformed, err := Transform(fm.Transform, valStr, fm.TransformArg)
		if err != nil {
			continue
		}

		switch fm.TargetField {
		case "report_type":
			report.ReportType = transformed
		case "report_no":
			report.ReportNo = sql.NullString{String: transformed, Valid: true}
		case "hospital_code":
			report.HospitalCode = sql.NullString{String: transformed, Valid: true}
		case "report_date":
			if t, err := time.Parse("2006-01-02", transformed); err == nil {
				report.ReportDate = sql.NullTime{Time: t, Valid: true}
			}
		case "summary":
			report.Summary = sql.NullString{String: transformed, Valid: true}
		case "abnormal_flags":
			report.AbnormalFlags = sql.NullString{String: transformed, Valid: true}
		case "source":
			report.Source = transformed
		}
	}
	return report, nil
}

func (m *SyncManager) writeToDeadLetter(source, rawData, errReason string) {
	entry := &model.DeadLetterEntry{
		Source:      source,
		RawData:     rawData,
		ErrorReason: errReason,
		Status:      "pending",
	}
	_ = store.NewDeadLetterStore(m.db).Insert(entry)
}
