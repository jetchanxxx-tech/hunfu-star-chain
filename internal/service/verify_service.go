package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/huifu/star-chain/internal/model"
	"github.com/huifu/star-chain/internal/store"
)

type VerifyService struct {
	store  *store.VerificationStore
	family *store.FamilyStore
	db     *sql.DB
	secret string
}

func NewVerifyService(db *sql.DB, secret string) *VerifyService {
	return &VerifyService{store: store.NewVerificationStore(db), family: store.NewFamilyStore(db), db: db, secret: secret}
}

type QRPayload struct {
	MemberUUID    string `json:"member_uuid"`
	EntitlementID int64  `json:"entitlement_id"`
	BenefitType   string `json:"benefit_type"`
	Timestamp     int64  `json:"timestamp"`
	Nonce         string `json:"nonce"`
	Signature     string `json:"signature"`
}

func (s *VerifyService) GenerateQRPayload(memberUUID string, entitlementID int64, benefitType string) (*QRPayload, error) {
	nonce := uuid.New().String()
	ts := time.Now().Unix()
	p := &QRPayload{
		MemberUUID:    memberUUID,
		EntitlementID: entitlementID,
		BenefitType:   benefitType,
		Timestamp:     ts,
		Nonce:         nonce,
	}
	p.Signature = p.ComputeSignature(s.secret)
	return p, nil
}

func (p *QRPayload) ComputeSignature(secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%s:%d:%d:%s:%s", p.MemberUUID, p.EntitlementID, p.Timestamp, p.BenefitType, p.Nonce)))
	return hex.EncodeToString(mac.Sum(nil))
}

func (p *QRPayload) Verify(secret string) bool {
	expected := p.ComputeSignature(secret)
	return hmac.Equal([]byte(p.Signature), []byte(expected))
}

func (p *QRPayload) IsExpired() bool {
	return time.Now().Unix()-p.Timestamp > 60
}

func (s *VerifyService) Verify(payload *QRPayload, stewardID int64) (*model.VerificationRecord, error) {
	if !payload.Verify(s.secret) {
		return nil, fmt.Errorf("invalid qr signature")
	}
	if payload.IsExpired() {
		return nil, fmt.Errorf("qr code expired")
	}
	used, err := s.store.QrNonceExists(payload.Nonce)
	if err != nil {
		return nil, err
	}
	if used {
		return nil, fmt.Errorf("qr code already used")
	}

	member, err := s.family.FindMemberByUUID(payload.MemberUUID)
	if err != nil {
		return nil, fmt.Errorf("member not found")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Atomic consume
	if err := s.store.ConsumeEntitlement(tx, payload.EntitlementID, 1); err != nil {
		return nil, fmt.Errorf("entitlement consumed: %w", err)
	}

	r := &model.VerificationRecord{
		EntitlementID: payload.EntitlementID,
		MemberID:      member.ID,
		StewardID:     sql.NullInt64{Int64: stewardID, Valid: stewardID > 0},
		QRNonce:       payload.Nonce,
		BenefitType:   payload.BenefitType,
		VerifyCount:   1,
		Status:        "success",
	}
	if err := s.store.InsertRecord(tx, r); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *VerifyService) ListRecords(memberID int64, offset, limit int) ([]model.VerificationRecord, error) {
	return s.store.ListByMember(memberID, offset, limit)
}

// DecodePayload parses a JSON string into QRPayload
func DecodePayload(data string) (*QRPayload, error) {
	var p QRPayload
	if err := json.Unmarshal([]byte(data), &p); err != nil {
		return nil, fmt.Errorf("decode qr payload: %w", err)
	}
	return &p, nil
}
