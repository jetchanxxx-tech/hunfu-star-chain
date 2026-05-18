package store

import (
	"database/sql"

	"github.com/huifu/star-chain/internal/model"
)

type TimelineStore struct{ db *sql.DB }

func NewTimelineStore(db *sql.DB) *TimelineStore { return &TimelineStore{db: db} }

func (s *TimelineStore) ListByMember(memberID int64, limit int) ([]model.TimelineEvent, error) {
	rows, err := s.db.Query(
		`SELECT id, member_id, event_type, event_date, event_data, source, created_at
		 FROM timeline_events WHERE member_id = ? ORDER BY event_date DESC LIMIT ?`,
		memberID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.TimelineEvent
	for rows.Next() {
		var e model.TimelineEvent
		if err := rows.Scan(&e.ID, &e.MemberID, &e.EventType, &e.EventDate, &e.EventData, &e.Source, &e.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

func (s *TimelineStore) Create(e *model.TimelineEvent) error {
	_, err := s.db.Exec(
		`INSERT INTO timeline_events (member_id, event_type, event_date, event_data, source)
		 VALUES (?, ?, ?, ?, ?)`,
		e.MemberID, e.EventType, e.EventDate, e.EventData, e.Source,
	)
	return err
}
