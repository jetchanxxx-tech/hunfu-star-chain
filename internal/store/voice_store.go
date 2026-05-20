package store

import (
	"database/sql"

	"github.com/huifu/star-chain/internal/model"
)

type VoiceStore struct{ db *sql.DB }

func NewVoiceStore(db *sql.DB) *VoiceStore { return &VoiceStore{db: db} }

func (s *VoiceStore) InsertLog(log *model.VoiceCallLog) error {
	res, err := s.db.Exec(`INSERT INTO voice_call_logs
		(task_id, member_id, phone, call_type, provider, template_code, dialogue_text, call_status, duration_seconds, transcript, fail_reason, called_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		log.TaskID, log.MemberID, log.Phone, log.CallType, log.Provider,
		log.TemplateCode, log.DialogueText, log.CallStatus,
		log.DurationSeconds, log.Transcript, log.FailReason)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	log.ID = id
	return nil
}

func (s *VoiceStore) UpdateCallResult(callID, status string, duration int, transcript, failReason string) error {
	_, err := s.db.Exec(`UPDATE voice_call_logs SET call_status = ?, duration_seconds = ?, transcript = ?, fail_reason = ?, finished_at = NOW()
		WHERE id = ?`, status, duration, transcript, failReason, callID)
	return err
}

func (s *VoiceStore) ListByMember(memberID int64, offset, limit int) ([]model.VoiceCallLog, error) {
	rows, err := s.db.Query(`SELECT id, task_id, member_id, phone, call_type, provider, template_code, dialogue_text, call_status, duration_seconds, transcript, intent_tags, retry_count, fail_reason, called_at, answered_at, finished_at, created_at
		FROM voice_call_logs WHERE member_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, memberID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.VoiceCallLog
	for rows.Next() {
		var v model.VoiceCallLog
		if err := rows.Scan(&v.ID, &v.TaskID, &v.MemberID, &v.Phone, &v.CallType, &v.Provider,
			&v.TemplateCode, &v.DialogueText, &v.CallStatus, &v.DurationSeconds,
			&v.Transcript, &v.IntentTags, &v.RetryCount, &v.FailReason,
			&v.CalledAt, &v.AnsweredAt, &v.FinishedAt, &v.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (s *VoiceStore) CountTodayCalls(memberID int64) (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM voice_call_logs
		WHERE member_id = ? AND DATE(created_at) = CURDATE() AND call_type = 'auto'`, memberID).Scan(&count)
	return count, err
}

func (s *VoiceStore) FindTemplateByCode(code string) (*model.VoiceTemplate, error) {
	t := &model.VoiceTemplate{}
	err := s.db.QueryRow(`SELECT id, template_code, template_name, category, provider, provider_tpl_id, llm_prompt, max_retries, retry_interval, daily_limit, status, created_at, updated_at
		FROM voice_templates WHERE template_code = ?`, code).Scan(
		&t.ID, &t.TemplateCode, &t.TemplateName, &t.Category, &t.Provider,
		&t.ProviderTplID, &t.LLMPrompt, &t.MaxRetries, &t.RetryInterval,
		&t.DailyLimit, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *VoiceStore) ListTemplates() ([]model.VoiceTemplate, error) {
	rows, err := s.db.Query(`SELECT id, template_code, template_name, category, provider, provider_tpl_id, llm_prompt, max_retries, retry_interval, daily_limit, status, created_at, updated_at
		FROM voice_templates ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.VoiceTemplate
	for rows.Next() {
		var t model.VoiceTemplate
		if err := rows.Scan(&t.ID, &t.TemplateCode, &t.TemplateName, &t.Category, &t.Provider,
			&t.ProviderTplID, &t.LLMPrompt, &t.MaxRetries, &t.RetryInterval,
			&t.DailyLimit, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *VoiceStore) UpsertTemplate(t *model.VoiceTemplate) error {
	_, err := s.db.Exec(`INSERT INTO voice_templates (template_code, template_name, category, provider, provider_tpl_id, llm_prompt, max_retries, retry_interval, daily_limit, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE template_name=VALUES(template_name), category=VALUES(category),
		provider=VALUES(provider), provider_tpl_id=VALUES(provider_tpl_id), llm_prompt=VALUES(llm_prompt),
		max_retries=VALUES(max_retries), retry_interval=VALUES(retry_interval), daily_limit=VALUES(daily_limit), status=VALUES(status)`,
		t.TemplateCode, t.TemplateName, t.Category, t.Provider, t.ProviderTplID,
		t.LLMPrompt, t.MaxRetries, t.RetryInterval, t.DailyLimit, t.Status)
	return err
}
