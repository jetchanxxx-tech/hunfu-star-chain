package router

import (
	"database/sql"
	"net/http"

	"github.com/huifu/star-chain/internal/config"
	"github.com/huifu/star-chain/internal/handler"
	"github.com/huifu/star-chain/internal/middleware"
)

func New(cfg *config.Config, db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Health
	healthH := handler.NewHealth(db)
	mux.HandleFunc("GET /api/health", healthH.Check)

	// Auth (WeChat login)
	authH := handler.NewAuthHandler(cfg, db)
	mux.HandleFunc("POST /api/v1/auth/wx-login", authH.WxLogin)
	mux.HandleFunc("POST /api/v1/auth/bind-wechat", authH.BindWechat)

	// Families & Members
	familyH := handler.NewFamilyHandler(db)
	mux.HandleFunc("POST /api/v1/families", familyH.Create)
	mux.HandleFunc("GET /api/v1/families/{id}", familyH.Get)
	mux.HandleFunc("POST /api/v1/families/{id}/members", familyH.AddMember)

	// Timeline
	timelineH := handler.NewTimelineHandler(db)
	mux.HandleFunc("GET /api/v1/members/{id}/timeline", timelineH.ListByMember)

	// Reports
	reportH := handler.NewReportHandler(db)
	mux.HandleFunc("GET /api/v1/members/{id}/reports", reportH.ListByMember)

	// Service Packages
	packageH := handler.NewPackageHandler(db)
	mux.HandleFunc("GET /api/v1/packages", packageH.List)

	// AI
	aiH := handler.NewAIHandler(cfg, db)
	mux.HandleFunc("POST /api/v1/ai/chat", aiH.Chat)
	mux.HandleFunc("GET /api/v1/ai/faq", aiH.FAQSearch)

	var h http.Handler = mux
	h = middleware.CORS(h)
	h = middleware.Logger(h)
	h = middleware.Recovery(h)
	return h
}
