package service

import (
	"crypto/sha256"
	"database/sql"
	"fmt"

	"github.com/huifu/star-chain/internal/model"
	"github.com/huifu/star-chain/internal/store"

	"github.com/google/uuid"
)

type FamilyService struct {
	family   *store.FamilyStore
	timeline *store.TimelineStore
	db       *sql.DB
}

func NewFamilyService(db *sql.DB) *FamilyService {
	return &FamilyService{
		family:   store.NewFamilyStore(db),
		timeline: store.NewTimelineStore(db),
		db:       db,
	}
}

func (s *FamilyService) Create(req model.CreateFamilyRequest) (*model.CreateFamilyResponse, error) {
	phoneHash := hashString(req.PhoneHash)
	existing, _ := s.family.FindMemberByPhoneHash(phoneHash)
	if existing != nil {
		return nil, fmt.Errorf("member already exists")
	}

	memberUUID := uuid.New().String()
	familyUUID := uuid.New().String()

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	f := &model.Family{UUID: familyUUID, Name: sql.NullString{String: req.Name, Valid: req.Name != ""}}
	if err := s.family.Create(tx, f); err != nil {
		return nil, fmt.Errorf("create family: %w", err)
	}

	// fetch the auto-incremented ID
	fObj, err := s.family.FindByUUID(familyUUID)
	if err != nil {
		return nil, fmt.Errorf("find family: %w", err)
	}

	m := &model.FamilyMember{
		UUID:         memberUUID,
		FamilyID:     fObj.ID,
		Relation:     normalizeRelation(req.Relation, "self"),
		Nickname:     sql.NullString{String: req.Nickname, Valid: req.Nickname != ""},
		RealNameHash: phoneHash,
		PhoneHash:    phoneHash,
	}
	if err := s.family.CreateMember(tx, m); err != nil {
		return nil, fmt.Errorf("create member: %w", err)
	}

	// update primary_member_id
	if _, err := tx.Exec("UPDATE families SET primary_member_id = (SELECT id FROM family_members WHERE member_uuid = ?) WHERE family_uuid = ?", memberUUID, familyUUID); err != nil {
		return nil, fmt.Errorf("update primary: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return &model.CreateFamilyResponse{
		FamilyID:   familyUUID,
		MemberID:   memberUUID,
		MemberUUID: memberUUID,
	}, nil
}

func (s *FamilyService) AddMember(familyUUID string, req model.AddMemberRequest) (*model.MemberResponse, error) {
	f, err := s.family.FindByUUID(familyUUID)
	if err != nil {
		return nil, fmt.Errorf("family not found: %w", err)
	}
	phoneHash := hashString(req.PhoneHash)
	m := &model.FamilyMember{
		UUID:         uuid.New().String(),
		FamilyID:     f.ID,
		Relation:     normalizeRelation(req.Relation, "other"),
		Nickname:     sql.NullString{String: req.Nickname, Valid: req.Nickname != ""},
		RealNameHash: phoneHash,
		PhoneHash:    phoneHash,
	}
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()
	if err := s.family.CreateMember(tx, m); err != nil {
		return nil, fmt.Errorf("create member: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return &model.MemberResponse{
		UUID:     m.UUID,
		FamilyID: familyUUID,
		Nickname: req.Nickname,
		Relation: m.Relation,
		Status:   "active",
	}, nil
}

func (s *FamilyService) FindMemberByUUID(uuid string) (*model.FamilyMember, error) {
	return s.family.FindMemberByUUID(uuid)
}

func (s *FamilyService) GetFamily(familyUUID string) (*model.FamilyResponse, error) {
	f, err := s.family.FindByUUID(familyUUID)
	if err != nil {
		return nil, fmt.Errorf("family not found: %w", err)
	}
	members, err := s.family.ListMembers(f.ID)
	if err != nil {
		return nil, fmt.Errorf("list members: %w", err)
	}
	resp := &model.FamilyResponse{
		UUID: f.UUID,
		Name: f.Name.String,
	}
	for _, m := range members {
		resp.Members = append(resp.Members, model.MemberResponse{
			UUID:      m.UUID,
			FamilyID:  f.UUID,
			Nickname:  m.Nickname.String,
			Relation:  m.Relation,
			BirthDate: nullTimeStr(m.BirthDate),
			Gender:    int(nullInt(m.Gender)),
			Status:    m.Status,
		})
	}
	return resp, nil
}

func hashString(s string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
}

func normalizeRelation(s, defaultRel string) string {
	valid := map[string]bool{"self": true, "spouse": true, "child": true, "parent": true, "other": true}
	if valid[s] {
		return s
	}
	return defaultRel
}

func nullTimeStr(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format("2006-01-02")
	}
	return ""
}

func nullInt(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}
