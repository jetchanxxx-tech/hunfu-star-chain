package model

import (
	"database/sql"
	"time"
)

type TimelineNodeTemplate struct {
	ID           int64          `json:"id"`
	NodeCode     string         `json:"node_code"`
	NodeName     string         `json:"node_name"`
	Category     string         `json:"category"`
	DefaultStart sql.NullString `json:"default_start"`
	DefaultEnd   sql.NullString `json:"default_end"`
	ReminderDays int            `json:"reminder_days"`
	SortOrder    int            `json:"sort_order"`
	Description  sql.NullString `json:"description"`
	Status       string         `json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type HospitalNodeOverride struct {
	ID           int64          `json:"id"`
	HospitalCode string         `json:"hospital_code"`
	NodeCode     string         `json:"node_code"`
	StartOffset  sql.NullString `json:"start_offset"`
	EndOffset    sql.NullString `json:"end_offset"`
	ReminderDays sql.NullInt64  `json:"reminder_days"`
	IsEnabled    bool           `json:"is_enabled"`
	Description  sql.NullString `json:"description"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

// Member authorization
type MemberAuthorization struct {
	ID           int64          `json:"id"`
	GrantorID    int64          `json:"grantor_id"`
	GranteeID    int64          `json:"grantee_id"`
	AuthScope    sql.NullString `json:"auth_scope"` // JSON array
	Status       string         `json:"status"`
	RequestMsg   sql.NullString `json:"request_msg"`
	RejectReason sql.NullString `json:"reject_reason"`
	ValidUntil   sql.NullTime   `json:"valid_until"`
	RequestedAt  time.Time      `json:"requested_at"`
	RespondedAt  sql.NullTime   `json:"responded_at"`
	RevokedAt    sql.NullTime   `json:"revoked_at"`
}

type AuthorizationAuditLog struct {
	ID              int64          `json:"id"`
	AuthorizationID sql.NullInt64  `json:"authorization_id"`
	Action          string         `json:"action"`
	ActorID         int64          `json:"actor_id"`
	TargetID        int64          `json:"target_id"`
	Detail          sql.NullString `json:"detail"` // JSON
	CreatedAt       time.Time      `json:"created_at"`
}

// Verification
type VerificationRecord struct {
	ID            int64          `json:"id"`
	EntitlementID int64          `json:"entitlement_id"`
	MemberID      int64          `json:"member_id"`
	StewardID     sql.NullInt64  `json:"steward_id"`
	QRNonce       string         `json:"qr_nonce"`
	BenefitType   string         `json:"benefit_type"`
	VerifyCount   int            `json:"verify_count"`
	Status        string         `json:"status"`
	FailReason    sql.NullString `json:"fail_reason"`
	VerifiedAt    time.Time      `json:"verified_at"`
	ConfirmedAt   sql.NullTime   `json:"confirmed_at"`
}

// Dead letter queue
type DeadLetterEntry struct {
	ID          int64          `json:"id"`
	Source      string         `json:"source"`
	RawData     string         `json:"raw_data"`
	ErrorReason string         `json:"error_reason"`
	RetryCount  int            `json:"retry_count"`
	Status      string         `json:"status"`
	ResolvedAt  sql.NullTime   `json:"resolved_at"`
	CreatedAt   time.Time      `json:"created_at"`
}

type DataSyncLog struct {
	ID           int64          `json:"id"`
	Source       string         `json:"source"`
	SyncType     string         `json:"sync_type"`
	TotalCount   int            `json:"total_count"`
	SuccessCount int            `json:"success_count"`
	FailCount    int            `json:"fail_count"`
	StartedAt    sql.NullTime   `json:"started_at"`
	FinishedAt   sql.NullTime   `json:"finished_at"`
	CreatedAt    time.Time      `json:"created_at"`
}

// Voice call
type VoiceCallLog struct {
	ID              int64          `json:"id"`
	TaskID          sql.NullInt64  `json:"task_id"`
	MemberID        int64          `json:"member_id"`
	Phone           string         `json:"phone"`
	CallType        string         `json:"call_type"`
	Provider        string         `json:"provider"`
	TemplateCode    sql.NullString `json:"template_code"`
	DialogueText    sql.NullString `json:"dialogue_text"`
	CallStatus      string         `json:"call_status"`
	DurationSeconds int            `json:"duration_seconds"`
	Transcript      sql.NullString `json:"transcript"`
	IntentTags      sql.NullString `json:"intent_tags"`
	RetryCount      int            `json:"retry_count"`
	FailReason      sql.NullString `json:"fail_reason"`
	CalledAt        sql.NullTime   `json:"called_at"`
	AnsweredAt      sql.NullTime   `json:"answered_at"`
	FinishedAt      sql.NullTime   `json:"finished_at"`
	CreatedAt       time.Time      `json:"created_at"`
}

type VoiceTemplate struct {
	ID             int64          `json:"id"`
	TemplateCode   string         `json:"template_code"`
	TemplateName   string         `json:"template_name"`
	Category       string         `json:"category"`
	Provider       string         `json:"provider"`
	ProviderTplID  sql.NullString `json:"provider_tpl_id"`
	LLMPrompt      sql.NullString `json:"llm_prompt"`
	MaxRetries     int            `json:"max_retries"`
	RetryInterval  int            `json:"retry_interval"`
	DailyLimit     int            `json:"daily_limit"`
	Status         string         `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
