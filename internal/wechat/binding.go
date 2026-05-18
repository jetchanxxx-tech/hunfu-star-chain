package wechat

import (
	"database/sql"
	"fmt"
)

type BindingStore struct{ db *sql.DB }

func NewBindingStore(db *sql.DB) *BindingStore { return &BindingStore{db: db} }

func (s *BindingStore) BindOpenid(memberID int64, openidMP, unionID string) error {
	_, err := s.db.Exec(
		`INSERT INTO wechat_bindings (member_id, openid_mp, union_id)
		 VALUES (?, ?, ?)
		 ON DUPLICATE KEY UPDATE openid_mp = VALUES(openid_mp), union_id = COALESCE(VALUES(union_id), union_id)`,
		memberID, openidMP, nullIfEmpty(unionID),
	)
	if err != nil {
		return fmt.Errorf("bind openid: %w", err)
	}
	return nil
}

func (s *BindingStore) FindMemberByOpenid(openidMP string) (int64, error) {
	var memberID int64
	err := s.db.QueryRow(
		"SELECT member_id FROM wechat_bindings WHERE openid_mp = ?", openidMP,
	).Scan(&memberID)
	if err != nil {
		return 0, err
	}
	return memberID, nil
}

func (s *BindingStore) BindWework(memberID int64, openidWework string) error {
	_, err := s.db.Exec(
		`UPDATE wechat_bindings SET openid_wework = ? WHERE member_id = ?`,
		openidWework, memberID,
	)
	return err
}

func nullIfEmpty(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
