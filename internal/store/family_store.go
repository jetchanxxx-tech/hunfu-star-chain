package store

import (
	"database/sql"
	"fmt"

	"github.com/huifu/star-chain/internal/model"
)

type FamilyStore struct{ db *sql.DB }

func NewFamilyStore(db *sql.DB) *FamilyStore { return &FamilyStore{db: db} }

func (s *FamilyStore) Create(tx *sql.Tx, f *model.Family) error {
	_, err := tx.Exec(
		`INSERT INTO families (family_uuid, name, primary_member_id) VALUES (?, ?, ?)`,
		f.UUID, f.Name, f.PrimaryMemberID,
	)
	return err
}

func (s *FamilyStore) CreateMember(tx *sql.Tx, m *model.FamilyMember) error {
	_, err := tx.Exec(
		`INSERT INTO family_members (member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.UUID, m.FamilyID, m.Relation, m.Nickname, m.RealNameHash, m.PhoneHash,
		m.Gender, m.BirthDate, m.HealthCard,
	)
	return err
}

func (s *FamilyStore) FindByUUID(uuid string) (*model.Family, error) {
	var f model.Family
	err := s.db.QueryRow(
		`SELECT id, family_uuid, name, primary_member_id, created_at, updated_at FROM families WHERE family_uuid = ?`, uuid,
	).Scan(&f.ID, &f.UUID, &f.Name, &f.PrimaryMemberID, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *FamilyStore) FindMemberByUUID(uuid string) (*model.FamilyMember, error) {
	var m model.FamilyMember
	err := s.db.QueryRow(
		`SELECT id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash,
		        gender, birth_date, health_card, avatar_url, status, created_at, updated_at
		 FROM family_members WHERE member_uuid = ?`, uuid,
	).Scan(&m.ID, &m.UUID, &m.FamilyID, &m.Relation, &m.Nickname, &m.RealNameHash, &m.PhoneHash,
		&m.Gender, &m.BirthDate, &m.HealthCard, &m.AvatarURL, &m.Status, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *FamilyStore) FindMemberByPhoneHash(hash string) (*model.FamilyMember, error) {
	var m model.FamilyMember
	err := s.db.QueryRow(
		`SELECT id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash,
		        gender, birth_date, health_card, avatar_url, status, created_at, updated_at
		 FROM family_members WHERE phone_hash = ? AND status = 'active'`, hash,
	).Scan(&m.ID, &m.UUID, &m.FamilyID, &m.Relation, &m.Nickname, &m.RealNameHash, &m.PhoneHash,
		&m.Gender, &m.BirthDate, &m.HealthCard, &m.AvatarURL, &m.Status, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (s *FamilyStore) ListMembers(familyID int64) ([]model.FamilyMember, error) {
	rows, err := s.db.Query(
		`SELECT id, member_uuid, family_id, relation, nickname,
		        gender, birth_date, avatar_url, status, created_at, updated_at
		 FROM family_members WHERE family_id = ? AND status = 'active'`, familyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []model.FamilyMember
	for rows.Next() {
		var m model.FamilyMember
		if err := rows.Scan(&m.ID, &m.UUID, &m.FamilyID, &m.Relation, &m.Nickname,
			&m.Gender, &m.BirthDate, &m.AvatarURL, &m.Status, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		members = append(members, m)
	}
	return members, rows.Err()
}
