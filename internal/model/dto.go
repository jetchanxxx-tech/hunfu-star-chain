package model

type CreateFamilyRequest struct {
	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	PhoneHash string `json:"phone_hash"`
	Relation  string `json:"relation"` // defaults to 'self'
}

type CreateFamilyResponse struct {
	FamilyID   string `json:"family_id"`
	MemberID   string `json:"member_id"`
	MemberUUID string `json:"member_uuid"`
}

type AddMemberRequest struct {
	Nickname  string `json:"nickname"`
	Relation  string `json:"relation"`
	PhoneHash string `json:"phone_hash"`
	BirthDate string `json:"birth_date,omitempty"`
	Gender    int    `json:"gender,omitempty"`
}

type MemberResponse struct {
	UUID      string `json:"member_uuid"`
	FamilyID  string `json:"family_id"`
	Nickname  string `json:"nickname"`
	Relation  string `json:"relation"`
	BirthDate string `json:"birth_date,omitempty"`
	Gender    int    `json:"gender,omitempty"`
	Status    string `json:"status"`
}

type FamilyResponse struct {
	UUID    string           `json:"family_uuid"`
	Name    string           `json:"name"`
	Members []MemberResponse `json:"members"`
}

type TimelineEventResponse struct {
	ID        int64  `json:"id"`
	EventType string `json:"event_type"`
	EventDate string `json:"event_date"`
	EventData string `json:"event_data,omitempty"`
	Source    string `json:"source"`
}

type ReportResponse struct {
	ID            int64  `json:"id"`
	ReportType    string `json:"report_type"`
	HospitalCode  string `json:"hospital_code,omitempty"`
	ReportDate    string `json:"report_date,omitempty"`
	Summary       string `json:"summary,omitempty"`
	AbnormalFlags string `json:"abnormal_flags,omitempty"`
	Source        string `json:"source"`
}

type PackageResponse struct {
	UUID        string  `json:"package_uuid"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Level       string  `json:"level"`
	Price       float64 `json:"price"`
	CoverImage  string  `json:"cover_image,omitempty"`
	Benefits    string  `json:"benefits,omitempty"`
	Status      string  `json:"status"`
}
