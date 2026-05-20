package service

import (
	"database/sql"
	"fmt"

	"github.com/huifu/star-chain/internal/model"
	"github.com/huifu/star-chain/internal/store"
)

type NodeService struct {
	store  *store.TimelineNodeStore
	family *store.FamilyStore
	db     *sql.DB
}

func NewNodeService(db *sql.DB) *NodeService {
	return &NodeService{store: store.NewTimelineNodeStore(db), family: store.NewFamilyStore(db), db: db}
}

func (s *NodeService) ListTemplates() ([]model.TimelineNodeTemplate, error) {
	return s.store.ListTemplates()
}

func (s *NodeService) UpsertTemplate(t *model.TimelineNodeTemplate) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()
	if err := s.store.UpsertTemplate(tx, t); err != nil {
		return fmt.Errorf("upsert template: %w", err)
	}
	return tx.Commit()
}

func (s *NodeService) UpdateTemplateStatus(code, status string) error {
	if status != "enabled" && status != "disabled" {
		return fmt.Errorf("invalid status: %s", status)
	}
	if _, err := s.store.FindTemplateByCode(code); err != nil {
		return fmt.Errorf("template not found: %s", code)
	}
	return s.store.UpdateTemplateStatus(code, status)
}

func (s *NodeService) ListOverrides(hospitalCode string) ([]model.HospitalNodeOverride, error) {
	return s.store.ListOverrides(hospitalCode)
}

func (s *NodeService) UpsertOverride(o *model.HospitalNodeOverride) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()
	if err := s.store.UpsertOverride(tx, o); err != nil {
		return fmt.Errorf("upsert override: %w", err)
	}
	return tx.Commit()
}

func (s *NodeService) DeleteOverride(hospitalCode, nodeCode string) error {
	return s.store.DeleteOverride(hospitalCode, nodeCode)
}
