package store

import (
	"database/sql"
	"time"

	"github.com/huifu/star-chain/internal/model"
)

type TimelineNodeStore struct{ db *sql.DB }

func NewTimelineNodeStore(db *sql.DB) *TimelineNodeStore { return &TimelineNodeStore{db: db} }

func (s *TimelineNodeStore) ListTemplates() ([]model.TimelineNodeTemplate, error) {
	rows, err := s.db.Query(`SELECT id, node_code, node_name, category, default_start, default_end,
		reminder_days, sort_order, description, status, created_at, updated_at
		FROM timeline_node_templates ORDER BY sort_order`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.TimelineNodeTemplate
	for rows.Next() {
		var t model.TimelineNodeTemplate
		if err := rows.Scan(&t.ID, &t.NodeCode, &t.NodeName, &t.Category,
			&t.DefaultStart, &t.DefaultEnd, &t.ReminderDays, &t.SortOrder,
			&t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *TimelineNodeStore) FindTemplateByCode(code string) (*model.TimelineNodeTemplate, error) {
	t := &model.TimelineNodeTemplate{}
	err := s.db.QueryRow(`SELECT id, node_code, node_name, category, default_start, default_end,
		reminder_days, sort_order, description, status, created_at, updated_at
		FROM timeline_node_templates WHERE node_code = ?`, code).Scan(
		&t.ID, &t.NodeCode, &t.NodeName, &t.Category,
		&t.DefaultStart, &t.DefaultEnd, &t.ReminderDays, &t.SortOrder,
		&t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *TimelineNodeStore) UpsertTemplate(tx *sql.Tx, t *model.TimelineNodeTemplate) error {
	_, err := tx.Exec(`INSERT INTO timeline_node_templates
		(node_code, node_name, category, default_start, default_end, reminder_days, sort_order, description, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE node_name=VALUES(node_name), category=VALUES(category),
		default_start=VALUES(default_start), default_end=VALUES(default_end),
		reminder_days=VALUES(reminder_days), sort_order=VALUES(sort_order),
		description=VALUES(description), status=VALUES(status)`,
		t.NodeCode, t.NodeName, t.Category, t.DefaultStart, t.DefaultEnd,
		t.ReminderDays, t.SortOrder, t.Description, t.Status)
	return err
}

func (s *TimelineNodeStore) UpdateTemplateStatus(code, status string) error {
	_, err := s.db.Exec(`UPDATE timeline_node_templates SET status = ? WHERE node_code = ?`, status, code)
	return err
}

func (s *TimelineNodeStore) ListOverrides(hospitalCode string) ([]model.HospitalNodeOverride, error) {
	rows, err := s.db.Query(`SELECT id, hospital_code, node_code, start_offset, end_offset,
		reminder_days, is_enabled, description, created_at, updated_at
		FROM hospital_node_overrides WHERE hospital_code = ?`, hospitalCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.HospitalNodeOverride
	for rows.Next() {
		var o model.HospitalNodeOverride
		if err := rows.Scan(&o.ID, &o.HospitalCode, &o.NodeCode,
			&o.StartOffset, &o.EndOffset, &o.ReminderDays,
			&o.IsEnabled, &o.Description, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func (s *TimelineNodeStore) UpsertOverride(tx *sql.Tx, o *model.HospitalNodeOverride) error {
	_, err := tx.Exec(`INSERT INTO hospital_node_overrides
		(hospital_code, node_code, start_offset, end_offset, reminder_days, is_enabled, description)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE start_offset=VALUES(start_offset), end_offset=VALUES(end_offset),
		reminder_days=VALUES(reminder_days), is_enabled=VALUES(is_enabled), description=VALUES(description)`,
		o.HospitalCode, o.NodeCode, o.StartOffset, o.EndOffset, o.ReminderDays, o.IsEnabled, o.Description)
	return err
}

func (s *TimelineNodeStore) DeleteOverride(hospitalCode, nodeCode string) error {
	_, err := s.db.Exec(`DELETE FROM hospital_node_overrides WHERE hospital_code = ? AND node_code = ?`, hospitalCode, nodeCode)
	return err
}

// --- Authorization store ---

type AuthStore struct{ db *sql.DB }

func NewAuthStore(db *sql.DB) *AuthStore { return &AuthStore{db: db} }

func (s *AuthStore) CreateAuthorization(tx *sql.Tx, a *model.MemberAuthorization) error {
	res, err := tx.Exec(`INSERT INTO member_authorizations
		(grantor_id, grantee_id, auth_scope, status, request_msg, valid_until)
		VALUES (?, ?, ?, ?, ?, ?)`,
		a.GrantorID, a.GranteeID, a.AuthScope, a.Status, a.RequestMsg, a.ValidUntil)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	a.ID = id
	return nil
}

func (s *AuthStore) FindPendingAuth(grantorID, granteeID int64) (*model.MemberAuthorization, error) {
	a := &model.MemberAuthorization{}
	err := s.db.QueryRow(`SELECT id, grantor_id, grantee_id, auth_scope, status, request_msg,
		reject_reason, valid_until, requested_at, responded_at, revoked_at
		FROM member_authorizations
		WHERE grantor_id = ? AND grantee_id = ? AND status = 'pending'
		ORDER BY requested_at DESC LIMIT 1`, grantorID, granteeID).Scan(
		&a.ID, &a.GrantorID, &a.GranteeID, &a.AuthScope, &a.Status, &a.RequestMsg,
		&a.RejectReason, &a.ValidUntil, &a.RequestedAt, &a.RespondedAt, &a.RevokedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AuthStore) FindAuthorizationByID(id int64) (*model.MemberAuthorization, error) {
	a := &model.MemberAuthorization{}
	err := s.db.QueryRow(`SELECT id, grantor_id, grantee_id, auth_scope, status, request_msg,
		reject_reason, valid_until, requested_at, responded_at, revoked_at
		FROM member_authorizations WHERE id = ?`, id).Scan(
		&a.ID, &a.GrantorID, &a.GranteeID, &a.AuthScope, &a.Status, &a.RequestMsg,
		&a.RejectReason, &a.ValidUntil, &a.RequestedAt, &a.RespondedAt, &a.RevokedAt)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AuthStore) UpdateAuthStatus(tx *sql.Tx, id int64, status string) error {
	now := time.Now()
	var setResponded, setRevoked string
	if status == "active" || status == "rejected" {
		setResponded = ", responded_at = ?"
	}
	if status == "revoked" {
		setRevoked = ", revoked_at = ?"
	}
	query := `UPDATE member_authorizations SET status = ?` + setResponded + setRevoked + ` WHERE id = ?`
	var args []interface{}
	args = append(args, status)
	if setResponded != "" {
		args = append(args, now)
	}
	if setRevoked != "" {
		args = append(args, now)
	}
	args = append(args, id)
	_, err := tx.Exec(query, args...)
	return err
}

func (s *AuthStore) ListAuthorizationsByGrantor(grantorID int64) ([]model.MemberAuthorization, error) {
	rows, err := s.db.Query(`SELECT id, grantor_id, grantee_id, auth_scope, status, request_msg,
		reject_reason, valid_until, requested_at, responded_at, revoked_at
		FROM member_authorizations WHERE grantor_id = ? ORDER BY requested_at DESC`, grantorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.MemberAuthorization
	for rows.Next() {
		var a model.MemberAuthorization
		if err := rows.Scan(&a.ID, &a.GrantorID, &a.GranteeID, &a.AuthScope, &a.Status,
			&a.RequestMsg, &a.RejectReason, &a.ValidUntil, &a.RequestedAt, &a.RespondedAt, &a.RevokedAt); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *AuthStore) CheckActiveAuth(grantorID, granteeID int64) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM member_authorizations
		WHERE grantor_id = ? AND grantee_id = ? AND status = 'active'
		AND (valid_until IS NULL OR valid_until >= CURDATE())`, grantorID, granteeID).Scan(&count)
	return count > 0, err
}

func (s *AuthStore) InsertAuditLog(tx *sql.Tx, log *model.AuthorizationAuditLog) error {
	_, err := tx.Exec(`INSERT INTO authorization_audit_log (authorization_id, action, actor_id, target_id, detail)
		VALUES (?, ?, ?, ?, ?)`, log.AuthorizationID, log.Action, log.ActorID, log.TargetID, log.Detail)
	return err
}

func (s *AuthStore) ListAuditLogs(offset, limit int) ([]model.AuthorizationAuditLog, error) {
	rows, err := s.db.Query(`SELECT id, authorization_id, action, actor_id, target_id, detail, created_at
		FROM authorization_audit_log ORDER BY created_at DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.AuthorizationAuditLog
	for rows.Next() {
		var l model.AuthorizationAuditLog
		if err := rows.Scan(&l.ID, &l.AuthorizationID, &l.Action, &l.ActorID, &l.TargetID, &l.Detail, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (s *AuthStore) ExpireAuthorizations() (int64, error) {
	res, err := s.db.Exec(`UPDATE member_authorizations SET status = 'expired'
		WHERE status = 'active' AND valid_until IS NOT NULL AND valid_until < CURDATE()`)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
