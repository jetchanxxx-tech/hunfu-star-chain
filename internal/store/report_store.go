package store

import (
	"database/sql"

	"github.com/huifu/star-chain/internal/model"
)

type ReportStore struct{ db *sql.DB }

func NewReportStore(db *sql.DB) *ReportStore { return &ReportStore{db: db} }

func (s *ReportStore) Insert(r *model.HealthReport) error {
	res, err := s.db.Exec(`INSERT INTO health_reports (member_id, report_type, hospital_code, report_no, summary, abnormal_flags, report_date, file_url, source)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.MemberID, r.ReportType, r.HospitalCode, r.ReportNo, r.Summary, r.AbnormalFlags, r.ReportDate, r.FileURL, r.Source)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	r.ID = id
	return nil
}

func (s *ReportStore) FindByReportNo(reportNo string) (*model.HealthReport, error) {
	r := &model.HealthReport{}
	err := s.db.QueryRow(`SELECT id, member_id, report_type, hospital_code, report_no,
		summary, abnormal_flags, report_date, file_url, source, created_at
		FROM health_reports WHERE report_no = ?`, reportNo).Scan(
		&r.ID, &r.MemberID, &r.ReportType, &r.HospitalCode, &r.ReportNo,
		&r.Summary, &r.AbnormalFlags, &r.ReportDate, &r.FileURL, &r.Source, &r.CreatedAt)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *ReportStore) ListByMember(memberID int64, limit, offset int) ([]model.HealthReport, error) {
	rows, err := s.db.Query(
		`SELECT id, member_id, report_type, hospital_code, report_no,
		        summary, abnormal_flags, report_date, file_url, source, created_at
		 FROM health_reports WHERE member_id = ? ORDER BY report_date DESC LIMIT ? OFFSET ?`,
		memberID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []model.HealthReport
	for rows.Next() {
		var r model.HealthReport
		if err := rows.Scan(&r.ID, &r.MemberID, &r.ReportType, &r.HospitalCode, &r.ReportNo,
			&r.Summary, &r.AbnormalFlags, &r.ReportDate, &r.FileURL, &r.Source, &r.CreatedAt); err != nil {
			return nil, err
		}
		reports = append(reports, r)
	}
	return reports, rows.Err()
}
