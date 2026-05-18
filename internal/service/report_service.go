package service

import (
	"database/sql"

	"github.com/huifu/star-chain/internal/model"
	"github.com/huifu/star-chain/internal/store"
)

type ReportService struct {
	report  *store.ReportStore
	family  *store.FamilyStore
}

func NewReportService(db *sql.DB) *ReportService {
	return &ReportService{
		report: store.NewReportStore(db),
		family: store.NewFamilyStore(db),
	}
}

func (s *ReportService) ListByMember(memberUUID string, limit, offset int) ([]model.ReportResponse, error) {
	member, err := s.family.FindMemberByUUID(memberUUID)
	if err != nil {
		return nil, err
	}
	reports, err := s.report.ListByMember(member.ID, limit, offset)
	if err != nil {
		return nil, err
	}
	var out []model.ReportResponse
	for _, r := range reports {
		out = append(out, model.ReportResponse{
			ID:            r.ID,
			ReportType:    r.ReportType,
			HospitalCode:  r.HospitalCode.String,
			ReportDate:    nullTimeStr(r.ReportDate),
			Summary:       r.Summary.String,
			AbnormalFlags: r.AbnormalFlags.String,
			Source:        r.Source,
		})
	}
	return out, nil
}
