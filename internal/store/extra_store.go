package store

import (
	"database/sql"
	"time"

	"github.com/huifu/star-chain/internal/model"
)

// --- Followup task store ---

type FollowupStore struct{ db *sql.DB }

func NewFollowupStore(db *sql.DB) *FollowupStore { return &FollowupStore{db: db} }

func (s *FollowupStore) Create(tx *sql.Tx, t *FollowupTask) error {
	res, err := tx.Exec(`INSERT INTO followup_tasks (member_id, trigger_type, trigger_value, title, status, assigned_to, due_date, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		t.MemberID, t.TriggerType, t.TriggerValue, t.Title, t.Status, t.AssignedTo, t.DueDate, t.Notes)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	t.ID = id
	return nil
}

type FollowupTask struct {
	ID           int64          `json:"id"`
	MemberID     int64          `json:"member_id"`
	TriggerType  string         `json:"trigger_type"`
	TriggerValue sql.NullString `json:"trigger_value"`
	Title        string         `json:"title"`
	Status       string         `json:"status"`
	AssignedTo   sql.NullInt64  `json:"assigned_to"`
	DueDate      sql.NullTime   `json:"due_date"`
	CompletedAt  sql.NullTime   `json:"completed_at"`
	Notes        sql.NullString `json:"notes"`
	CreatedAt    time.Time      `json:"created_at"`
}

func (s *FollowupStore) List(filters map[string]interface{}, offset, limit int) ([]FollowupTask, error) {
	query := `SELECT id, member_id, trigger_type, trigger_value, title, status, assigned_to, due_date, completed_at, notes, created_at
		FROM followup_tasks WHERE 1=1`
	args := []interface{}{}
	if v, ok := filters["status"]; ok && v != "" {
		query += " AND status = ?"
		args = append(args, v)
	}
	if v, ok := filters["trigger_type"]; ok && v != "" {
		query += " AND trigger_type = ?"
		args = append(args, v)
	}
	if v, ok := filters["assigned_to"]; ok && v.(int64) > 0 {
		query += " AND assigned_to = ?"
		args = append(args, v)
	}
	if v, ok := filters["member_id"]; ok && v.(int64) > 0 {
		query += " AND member_id = ?"
		args = append(args, v)
	}
	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []FollowupTask
	for rows.Next() {
		var t FollowupTask
		if err := rows.Scan(&t.ID, &t.MemberID, &t.TriggerType, &t.TriggerValue, &t.Title,
			&t.Status, &t.AssignedTo, &t.DueDate, &t.CompletedAt, &t.Notes, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *FollowupStore) FindByID(id int64) (*FollowupTask, error) {
	t := &FollowupTask{}
	err := s.db.QueryRow(`SELECT id, member_id, trigger_type, trigger_value, title, status, assigned_to, due_date, completed_at, notes, created_at
		FROM followup_tasks WHERE id = ?`, id).Scan(
		&t.ID, &t.MemberID, &t.TriggerType, &t.TriggerValue, &t.Title,
		&t.Status, &t.AssignedTo, &t.DueDate, &t.CompletedAt, &t.Notes, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *FollowupStore) Assign(tx *sql.Tx, id, stewardID int64) error {
	_, err := tx.Exec(`UPDATE followup_tasks SET assigned_to = ?, status = 'in_progress' WHERE id = ? AND status = 'pending'`, stewardID, id)
	return err
}

func (s *FollowupStore) Complete(tx *sql.Tx, id int64, notes string) error {
	_, err := tx.Exec(`UPDATE followup_tasks SET status = 'completed', completed_at = NOW(), notes = CONCAT(IFNULL(notes,''), '\n', ?) WHERE id = ?`, notes, id)
	return err
}

func (s *FollowupStore) Cancel(tx *sql.Tx, id int64, reason string) error {
	_, err := tx.Exec(`UPDATE followup_tasks SET status = 'cancelled', notes = CONCAT(IFNULL(notes,''), '\n取消原因: ', ?) WHERE id = ?`, reason, id)
	return err
}

func (s *FollowupStore) FindOverdue(dueDate time.Time) ([]FollowupTask, error) {
	rows, err := s.db.Query(`SELECT id, member_id, trigger_type, trigger_value, title, status, assigned_to, due_date, completed_at, notes, created_at
		FROM followup_tasks WHERE status IN ('pending','in_progress') AND due_date < ?`, dueDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []FollowupTask
	for rows.Next() {
		var t FollowupTask
		if err := rows.Scan(&t.ID, &t.MemberID, &t.TriggerType, &t.TriggerValue, &t.Title,
			&t.Status, &t.AssignedTo, &t.DueDate, &t.CompletedAt, &t.Notes, &t.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *FollowupStore) Stats() (map[string]int, error) {
	stats := map[string]int{}
	var n int
	for _, st := range []string{"pending", "in_progress"} {
		s.db.QueryRow(`SELECT COUNT(*) FROM followup_tasks WHERE status = ?`, st).Scan(&n)
		stats[st] = n
	}
	s.db.QueryRow(`SELECT COUNT(*) FROM followup_tasks WHERE status = 'completed' AND DATE(completed_at) = CURDATE()`).Scan(&n)
	stats["completed_today"] = n
	s.db.QueryRow(`SELECT COUNT(*) FROM followup_tasks WHERE status IN ('pending','in_progress') AND due_date < NOW()`).Scan(&n)
	stats["overdue"] = n
	return stats, nil
}

// --- Verification store ---

type VerificationStore struct{ db *sql.DB }

func NewVerificationStore(db *sql.DB) *VerificationStore { return &VerificationStore{db: db} }

func (s *VerificationStore) InsertRecord(tx *sql.Tx, r *model.VerificationRecord) error {
	res, err := tx.Exec(`INSERT INTO verification_records
		(entitlement_id, member_id, steward_id, qr_nonce, benefit_type, verify_count, status, fail_reason)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		r.EntitlementID, r.MemberID, r.StewardID, r.QRNonce, r.BenefitType, r.VerifyCount, r.Status, r.FailReason)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	r.ID = id
	return nil
}

func (s *VerificationStore) QrNonceExists(nonce string) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM verification_records WHERE qr_nonce = ?`, nonce).Scan(&count)
	return count > 0, err
}

func (s *VerificationStore) ListByMember(memberID int64, offset, limit int) ([]model.VerificationRecord, error) {
	rows, err := s.db.Query(`SELECT id, entitlement_id, member_id, steward_id, qr_nonce, benefit_type, verify_count, status, fail_reason, verified_at, confirmed_at
		FROM verification_records WHERE member_id = ? ORDER BY verified_at DESC LIMIT ? OFFSET ?`, memberID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.VerificationRecord
	for rows.Next() {
		var r model.VerificationRecord
		if err := rows.Scan(&r.ID, &r.EntitlementID, &r.MemberID, &r.StewardID, &r.QRNonce,
			&r.BenefitType, &r.VerifyCount, &r.Status, &r.FailReason, &r.VerifiedAt, &r.ConfirmedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *VerificationStore) ConsumeEntitlement(tx *sql.Tx, entitlementID int64, count int) error {
	_, err := tx.Exec(`UPDATE user_entitlements SET consumed = consumed + ? WHERE id = ? AND consumed + ? <= total`, count, entitlementID, count)
	return err
}
