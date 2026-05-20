package voice

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/huifu/star-chain/internal/model"
	"github.com/huifu/star-chain/internal/store"
)

type CallResult struct {
	CallID          string `json:"call_id"`
	Status          string `json:"status"` // completed/failed/rejected/no_answer
	DurationSeconds int    `json:"duration_seconds"`
	Transcript      string `json:"transcript"`
	FailReason      string `json:"fail_reason"`
}

type Provider interface {
	Name() string
	MakeCall(phone, dialogueText, templateID string, maxRetries int) (*CallResult, error)
	QueryCallStatus(callID string) (*CallResult, error)
	CallbackHandler(rawBody []byte) (*CallResult, error)
}

type Service struct {
	providers map[string]Provider
	defaultP  string
	db        *sql.DB
	store     *store.VoiceStore
	mu        sync.RWMutex
}

func NewService(db *sql.DB, defaultProvider string) *Service {
	return &Service{
		providers: make(map[string]Provider),
		defaultP:  defaultProvider,
		db:        db,
		store:     store.NewVoiceStore(db),
	}
}

func (s *Service) Register(p Provider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providers[p.Name()] = p
}

func (s *Service) getProvider(name string) Provider {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if name == "" {
		name = s.defaultP
	}
	return s.providers[name]
}

func (s *Service) MakeCall(memberID int64, phone, callType, provider, templateCode string) (*model.VoiceCallLog, error) {
	p := s.getProvider(provider)
	if p == nil {
		return nil, fmt.Errorf("voice provider not found: %s", provider)
	}

	var dialogueText string
	if templateCode != "" {
		tpl, err := s.store.FindTemplateByCode(templateCode)
		if err == nil && tpl.LLMPrompt.Valid {
			// In production, this would render via LLM with member context
			dialogueText = tpl.LLMPrompt.String
		}
	}

	log := &model.VoiceCallLog{
		MemberID:     memberID,
		Phone:        phone,
		CallType:     callType,
		Provider:     p.Name(),
		TemplateCode: sql.NullString{String: templateCode, Valid: templateCode != ""},
		DialogueText: sql.NullString{String: dialogueText, Valid: dialogueText != ""},
		CallStatus:   "pending",
	}

	result, err := p.MakeCall(phone, dialogueText, templateCode, 3)
	if err != nil {
		log.CallStatus = "failed"
		log.FailReason = sql.NullString{String: err.Error(), Valid: true}
	} else {
		log.CallStatus = result.Status
		log.DurationSeconds = result.DurationSeconds
		log.Transcript = sql.NullString{String: result.Transcript, Valid: result.Transcript != ""}
		if result.FailReason != "" {
			log.FailReason = sql.NullString{String: result.FailReason, Valid: true}
		}
	}

	if err := s.store.InsertLog(log); err != nil {
		return nil, fmt.Errorf("store call log: %w", err)
	}
	return log, nil
}

func (s *Service) ProcessCallback(providerName string, rawBody []byte) error {
	p := s.getProvider(providerName)
	if p == nil {
		return fmt.Errorf("provider not found: %s", providerName)
	}
	result, err := p.CallbackHandler(rawBody)
	if err != nil {
		return err
	}
	return s.store.UpdateCallResult(result.CallID, result.Status, result.DurationSeconds, result.Transcript, result.FailReason)
}

func (s *Service) ListTemplates() ([]model.VoiceTemplate, error) {
	return s.store.ListTemplates()
}

func (s *Service) UpsertTemplate(t *model.VoiceTemplate) error {
	return s.store.UpsertTemplate(t)
}

func (s *Service) CheckDailyLimit(memberID int64) (bool, error) {
	count, err := s.store.CountTodayCalls(memberID)
	if err != nil {
		return false, err
	}
	// Default limit 2 calls per day per member
	return count < 2, nil
}

func (s *Service) GetCallLogs(memberID int64, offset, limit int) ([]model.VoiceCallLog, error) {
	return s.store.ListByMember(memberID, offset, limit)
}

// GenerateDialogue generates AI dialogue using LLM with member context
func (s *Service) GenerateDialogue(providerName, templatePrompt string, context map[string]string) (string, error) {
	// Replace template variables with member context
	dialogue := templatePrompt
	for k, v := range context {
		dialogue = fmt.Sprintf(dialogue, k, v) // simplified; production would use proper templating
	}
	return dialogue, nil
}
