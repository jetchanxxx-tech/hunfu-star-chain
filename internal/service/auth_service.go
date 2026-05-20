package service

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/huifu/star-chain/internal/model"
	"github.com/huifu/star-chain/internal/store"
)

type AuthService struct {
	store  *store.AuthStore
	family *store.FamilyStore
	db     *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{store: store.NewAuthStore(db), family: store.NewFamilyStore(db), db: db}
}

func (s *AuthService) RequestAuthorization(grantorUUID, granteeUUID string, scopes []string, msg string) (*model.MemberAuthorization, error) {
	grantor, err := s.family.FindMemberByUUID(grantorUUID)
	if err != nil {
		return nil, fmt.Errorf("grantor not found: %w", err)
	}
	grantee, err := s.family.FindMemberByUUID(granteeUUID)
	if err != nil {
		return nil, fmt.Errorf("grantee not found: %w", err)
	}
	if grantor.ID == grantee.ID {
		return nil, fmt.Errorf("cannot authorize self")
	}

	// Check existing pending
	existing, _ := s.store.FindPendingAuth(grantor.ID, grantee.ID)
	if existing != nil {
		return nil, fmt.Errorf("pending authorization already exists (id=%d)", existing.ID)
	}

	scopeJSON, _ := json.Marshal(scopes)
	a := &model.MemberAuthorization{
		GrantorID: grantor.ID,
		GranteeID: grantee.ID,
		AuthScope: sql.NullString{String: string(scopeJSON), Valid: true},
		Status:    "pending",
		RequestMsg: sql.NullString{String: msg, Valid: msg != ""},
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.store.CreateAuthorization(tx, a); err != nil {
		return nil, err
	}

	// Audit log
	log := &model.AuthorizationAuditLog{
		AuthorizationID: sql.NullInt64{Int64: a.ID, Valid: true},
		Action:          "request",
		ActorID:         grantor.ID,
		TargetID:        grantee.ID,
		Detail:          a.AuthScope,
	}
	_ = s.store.InsertAuditLog(tx, log)

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AuthService) RespondToAuth(authID int64, granteeUUID, action string, rejectReason string) (*model.MemberAuthorization, error) {
	a, err := s.store.FindAuthorizationByID(authID)
	if err != nil {
		return nil, fmt.Errorf("authorization not found")
	}
	grantee, err := s.family.FindMemberByUUID(granteeUUID)
	if err != nil || grantee.ID != a.GranteeID {
		return nil, fmt.Errorf("not authorized to respond")
	}
	if a.Status != "pending" {
		return nil, fmt.Errorf("authorization is not pending")
	}

	var newStatus string
	switch action {
	case "approve":
		newStatus = "active"
	case "reject":
		newStatus = "rejected"
	default:
		return nil, fmt.Errorf("invalid action: %s", action)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.store.UpdateAuthStatus(tx, authID, newStatus); err != nil {
		return nil, err
	}

	log := &model.AuthorizationAuditLog{
		AuthorizationID: sql.NullInt64{Int64: authID, Valid: true},
		Action:          action,
		ActorID:         grantee.ID,
		TargetID:        a.GrantorID,
	}
	detail := map[string]string{}
	if rejectReason != "" {
		detail["reject_reason"] = rejectReason
	}
	b, _ := json.Marshal(detail)
	log.Detail = sql.NullString{String: string(b), Valid: true}
	_ = s.store.InsertAuditLog(tx, log)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	a.Status = newStatus
	return a, nil
}

func (s *AuthService) RevokeAuth(authID int64, grantorUUID string) error {
	a, err := s.store.FindAuthorizationByID(authID)
	if err != nil {
		return fmt.Errorf("authorization not found")
	}
	grantor, err := s.family.FindMemberByUUID(grantorUUID)
	if err != nil || grantor.ID != a.GrantorID {
		return fmt.Errorf("not authorized to revoke")
	}
	if a.Status != "active" {
		return fmt.Errorf("can only revoke active authorization")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.store.UpdateAuthStatus(tx, authID, "revoked"); err != nil {
		return err
	}

	log := &model.AuthorizationAuditLog{
		AuthorizationID: sql.NullInt64{Int64: authID, Valid: true},
		Action:          "revoke",
		ActorID:         grantor.ID,
		TargetID:        a.GranteeID,
	}
	_ = s.store.InsertAuditLog(tx, log)

	return tx.Commit()
}

func (s *AuthService) CheckAccess(memberUUID, targetUUID string) (bool, error) {
	member, err := s.family.FindMemberByUUID(memberUUID)
	if err != nil {
		return false, err
	}
	target, err := s.family.FindMemberByUUID(targetUUID)
	if err != nil {
		return false, err
	}
	if member.ID == target.ID {
		return true, nil
	}
	// Self or same family
	if member.FamilyID == target.FamilyID && target.Relation == "self" {
		return true, nil
	}
	// Check active authorization
	return s.store.CheckActiveAuth(target.ID, member.ID)
}

func (s *AuthService) ListMyAuthorizations(memberUUID string) ([]model.MemberAuthorization, error) {
	member, err := s.family.FindMemberByUUID(memberUUID)
	if err != nil {
		return nil, err
	}
	return s.store.ListAuthorizationsByGrantor(member.ID)
}

func (s *AuthService) ListAuditLogs(offset, limit int) ([]model.AuthorizationAuditLog, error) {
	return s.store.ListAuditLogs(offset, limit)
}

func (s *AuthService) ExpireAuthorizations() (int64, error) {
	return s.store.ExpireAuthorizations()
}
