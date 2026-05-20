package store

import (
	"database/sql"

	"github.com/huifu/star-chain/internal/model"
)

type DeadLetterStore struct{ db *sql.DB }

func NewDeadLetterStore(db *sql.DB) *DeadLetterStore { return &DeadLetterStore{db: db} }

func (s *DeadLetterStore) Insert(e *model.DeadLetterEntry) error {
	res, err := s.db.Exec(`INSERT INTO dead_letter_queue (source, raw_data, error_reason, status) VALUES (?, ?, ?, ?)`,
		e.Source, e.RawData, e.ErrorReason, e.Status)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	e.ID = id
	return nil
}

func (s *DeadLetterStore) List(source string, status string, offset, limit int) ([]model.DeadLetterEntry, error) {
	query := `SELECT id, source, raw_data, error_reason, retry_count, status, resolved_at, created_at FROM dead_letter_queue WHERE 1=1`
	args := []interface{}{}
	if source != "" {
		query += " AND source = ?"
		args = append(args, source)
	}
	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.DeadLetterEntry
	for rows.Next() {
		var e model.DeadLetterEntry
		if err := rows.Scan(&e.ID, &e.Source, &e.RawData, &e.ErrorReason, &e.RetryCount, &e.Status, &e.ResolvedAt, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *DeadLetterStore) Resolve(id int64) error {
	_, err := s.db.Exec(`UPDATE dead_letter_queue SET status = 'resolved', resolved_at = NOW() WHERE id = ?`, id)
	return err
}

func (s *DeadLetterStore) IncrementRetry(id int64) error {
	_, err := s.db.Exec(`UPDATE dead_letter_queue SET retry_count = retry_count + 1, status = 'retrying' WHERE id = ?`, id)
	return err
}
