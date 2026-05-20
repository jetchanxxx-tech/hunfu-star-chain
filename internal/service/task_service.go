package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/huifu/star-chain/internal/store"
)

type TaskService struct {
	store  *store.FollowupStore
	family *store.FamilyStore
	db     *sql.DB
}

func NewTaskService(db *sql.DB) *TaskService {
	return &TaskService{store: store.NewFollowupStore(db), family: store.NewFamilyStore(db), db: db}
}

func (s *TaskService) CreateTask(memberID int64, triggerType, triggerValue, title string, dueDate time.Time) (*store.FollowupTask, error) {
	t := &store.FollowupTask{
		MemberID:     memberID,
		TriggerType:  triggerType,
		TriggerValue: sql.NullString{String: triggerValue, Valid: triggerValue != ""},
		Title:        title,
		Status:       "pending",
		DueDate:      sql.NullTime{Time: dueDate, Valid: true},
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if err := s.store.Create(tx, t); err != nil {
		return nil, err
	}
	return t, tx.Commit()
}

func (s *TaskService) ListTasks(filters map[string]interface{}, offset, limit int) ([]store.FollowupTask, error) {
	return s.store.List(filters, offset, limit)
}

func (s *TaskService) GetTask(id int64) (*store.FollowupTask, error) {
	return s.store.FindByID(id)
}

func (s *TaskService) AssignTask(id, stewardID int64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.store.Assign(tx, id, stewardID); err != nil {
		return fmt.Errorf("assign: %w", err)
	}
	return tx.Commit()
}

func (s *TaskService) CompleteTask(id int64, notes string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.store.Complete(tx, id, notes); err != nil {
		return fmt.Errorf("complete: %w", err)
	}
	return tx.Commit()
}

func (s *TaskService) CancelTask(id int64, reason string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := s.store.Cancel(tx, id, reason); err != nil {
		return fmt.Errorf("cancel: %w", err)
	}
	return tx.Commit()
}

func (s *TaskService) CheckOverdue() ([]store.FollowupTask, error) {
	return s.store.FindOverdue(time.Now())
}

func (s *TaskService) Stats() (map[string]int, error) {
	return s.store.Stats()
}
