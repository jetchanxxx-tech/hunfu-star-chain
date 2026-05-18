package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/huifu/star-chain/internal/ai"
	"github.com/huifu/star-chain/internal/config"
)

type AIHandler struct {
	provider  ai.Provider
	faq       *ai.FAQMatcher
	emergency *ai.EmergencyDetector
	db        *sql.DB
}

func NewAIHandler(cfg *config.Config, db *sql.DB) *AIHandler {
	return &AIHandler{
		provider:  ai.NewProvider(cfg.AI),
		faq:       ai.NewFAQMatcher(db),
		emergency: ai.NewEmergencyDetector(cfg.AI.EmergencyKeywords),
		db:        db,
	}
}

type chatRequest struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
}

type chatResponse struct {
	Reply       string              `json:"reply"`
	Source      string              `json:"source"` // faq / ai
	Emergency   ai.EmergencyResult  `json:"emergency"`
	Emotion     string              `json:"emotion,omitempty"`
	SessionID   string              `json:"session_id,omitempty"`
}

func (h *AIHandler) Chat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.Message == "" {
		Error(w, http.StatusBadRequest, "message is required")
		return
	}

	emergency := h.emergency.Check(req.Message)

	var reply, source, emotion string

	// L1: try FAQ
	if entry, err := h.faq.Match(req.Message); err == nil && entry != nil {
		reply = entry.Answer
		source = "faq"
	} else {
		// L3: call LLM
		resp, err := h.provider.Chat(r.Context(), ai.ChatRequest{
			Messages: []ai.Message{
				{Role: "system", Content: "你是惠福星链的健康助手，帮助孕妈和家庭成员解答健康相关问题。回答需专业、温暖、简洁。"},
				{Role: "user", Content: req.Message},
			},
		})
		if err != nil {
			reply = "抱歉，AI服务暂时不可用，请稍后再试。紧急情况请拨打医院电话。"
			source = "fallback"
		} else if len(resp.Choices) > 0 {
			reply = resp.Choices[0].Message.Content
			source = "ai"
		}
	}

	// emotion detection
	_ = emotion // reserved for emotion detection via AI provider

	JSON(w, http.StatusOK, chatResponse{
		Reply:     reply,
		Source:    source,
		Emergency: emergency,
		SessionID: req.SessionID,
	})
}

func (h *AIHandler) FAQSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		Error(w, http.StatusBadRequest, "q parameter required")
		return
	}
	entry, err := h.faq.Match(query)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if entry == nil {
		JSON(w, http.StatusOK, map[string]any{"found": false})
		return
	}
	JSON(w, http.StatusOK, map[string]any{
		"found":    true,
		"question": entry.Question,
		"answer":   entry.Answer,
		"category": entry.Category,
	})
}
