package model

import (
	"database/sql"
	"time"
)

type Family struct {
	ID              int64          `json:"id"`
	UUID            string         `json:"uuid"`
	Name            sql.NullString `json:"name"`
	PrimaryMemberID int64          `json:"primary_member_id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type FamilyMember struct {
	ID             int64          `json:"id"`
	UUID           string         `json:"uuid"`
	FamilyID       int64          `json:"family_id"`
	Relation       string         `json:"relation"`
	Nickname       sql.NullString `json:"nickname"`
	RealNameHash   string         `json:"-"`
	PhoneHash      string         `json:"-"`
	EncryptedName  sql.NullString `json:"-"`
	EncryptedPhone sql.NullString `json:"-"`
	Gender         sql.NullInt64  `json:"gender"`
	BirthDate      sql.NullTime   `json:"birth_date"`
	HealthCard     sql.NullString `json:"health_card"`
	AvatarURL      sql.NullString `json:"avatar_url"`
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

type TimelineEvent struct {
	ID        int64          `json:"id"`
	MemberID  int64          `json:"member_id"`
	EventType string         `json:"event_type"`
	EventDate time.Time      `json:"event_date"`
	EventData sql.NullString `json:"event_data"` // JSON
	Source    string         `json:"source"`
	CreatedAt time.Time      `json:"created_at"`
}

type HealthReport struct {
	ID            int64          `json:"id"`
	MemberID      int64          `json:"member_id"`
	ReportType    string         `json:"report_type"`
	HospitalCode  sql.NullString `json:"hospital_code"`
	ReportNo      sql.NullString `json:"report_no"`
	Summary       sql.NullString `json:"summary"`       // JSON
	AbnormalFlags sql.NullString `json:"abnormal_flags"` // JSON
	ReportDate    sql.NullTime   `json:"report_date"`
	FileURL       sql.NullString `json:"file_url"`
	Source        string         `json:"source"`
	CreatedAt     time.Time      `json:"created_at"`
}

type ServicePackage struct {
	ID          int64          `json:"id"`
	UUID        string         `json:"uuid"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	Level       string         `json:"level"`
	Price       float64        `json:"price"`
	CoverImage  sql.NullString `json:"cover_image"`
	Benefits    sql.NullString `json:"benefits"` // JSON
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type UserEntitlement struct {
	ID          int64     `json:"id"`
	MemberID    int64     `json:"member_id"`
	PackageID   int64     `json:"package_id"`
	BenefitType string    `json:"benefit_type"`
	Total       int       `json:"total"`
	Consumed    int       `json:"consumed"`
	ValidUntil  time.Time `json:"valid_until"`
	CreatedAt   time.Time `json:"created_at"`
}
