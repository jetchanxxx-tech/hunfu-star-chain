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

	// Demo (公开，无需登录)
	demoH := handler.NewDemoHandler(db)
	mux.HandleFunc("GET /api/v1/demo/home", demoH.Home)

	// Auth (WeChat login)
	authH := handler.NewAuthHandler(cfg, db)
	mux.HandleFunc("POST /api/v1/auth/wx-login", authH.WxLogin)
	mux.HandleFunc("POST /api/v1/auth/bind-wechat", authH.BindWechat)

	// Admin
	adminH := handler.NewAdminHandler(db)
	mux.HandleFunc("POST /api/v1/admin/login", adminH.Login)
	mux.HandleFunc("GET /api/v1/admin/dashboard", adminH.Dashboard)
	mux.HandleFunc("GET /api/v1/admin/members", adminH.ListMembers)
	mux.HandleFunc("GET /api/v1/admin/packages", adminH.ListPackages)
	mux.HandleFunc("POST /api/v1/admin/packages", adminH.CreatePackage)
	mux.HandleFunc("PUT /api/v1/admin/packages/{id}", adminH.UpdatePackage)
	mux.HandleFunc("GET /api/v1/admin/users", adminH.ListUsers)

	// Families & Members
	familyH := handler.NewFamilyHandler(db)
	mux.HandleFunc("POST /api/v1/families", familyH.Create)
	mux.HandleFunc("GET /api/v1/families/{id}", familyH.Get)
	mux.HandleFunc("POST /api/v1/families/{id}/members", familyH.AddMember)

	// Timeline
	timelineH := handler.NewTimelineHandler(db)
	mux.HandleFunc("GET /api/v1/members/{id}/timeline", timelineH.ListByMember)

	// Timeline Node Templates (admin)
	nodeH := handler.NewNodeHandler(db)
	mux.HandleFunc("GET /api/v1/admin/node-templates", nodeH.ListTemplates)
	mux.HandleFunc("POST /api/v1/admin/node-templates", nodeH.UpsertTemplate)
	mux.HandleFunc("PATCH /api/v1/admin/node-templates/{code}/status", nodeH.UpdateTemplateStatus)
	mux.HandleFunc("GET /api/v1/admin/node-overrides/{hospitalCode}", nodeH.ListOverrides)
	mux.HandleFunc("POST /api/v1/admin/node-overrides", nodeH.UpsertOverride)
	mux.HandleFunc("DELETE /api/v1/admin/node-overrides/{hospitalCode}", nodeH.DeleteOverride)

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

	// Family Authorization
	authzH := handler.NewAuthzHandler(db)
	mux.HandleFunc("POST /api/v1/authorizations", authzH.RequestAuthorization)
	mux.HandleFunc("PATCH /api/v1/authorizations/{id}/respond", authzH.RespondToAuthorization)
	mux.HandleFunc("POST /api/v1/authorizations/{id}/revoke", authzH.RevokeAuthorization)
	mux.HandleFunc("GET /api/v1/authorizations", authzH.ListAuthorizations)
	mux.HandleFunc("GET /api/v1/authorizations/check", authzH.CheckAccess)
	mux.HandleFunc("GET /api/v1/admin/authorization-logs", authzH.ListAuditLogs)

	// Followup Tasks
	taskH := handler.NewTaskHandler(db)
	mux.HandleFunc("POST /api/v1/tasks", taskH.Create)
	mux.HandleFunc("GET /api/v1/tasks", taskH.List)
	mux.HandleFunc("GET /api/v1/tasks/{id}", taskH.Get)
	mux.HandleFunc("PATCH /api/v1/tasks/{id}/assign", taskH.Assign)
	mux.HandleFunc("PATCH /api/v1/tasks/{id}/complete", taskH.Complete)
	mux.HandleFunc("PATCH /api/v1/tasks/{id}/cancel", taskH.Cancel)
	mux.HandleFunc("GET /api/v1/admin/task-stats", taskH.Stats)

	// Verification
	verifySecret := cfg.QRVerifySecret
	if verifySecret == "" {
		verifySecret = "huifu-default-verify-secret"
	}
	verifyH := handler.NewVerifyHandler(db, verifySecret)
	mux.HandleFunc("POST /api/v1/verify/generate-qr", verifyH.GenerateQR)
	mux.HandleFunc("POST /api/v1/verify/consume", verifyH.Verify)
	mux.HandleFunc("GET /api/v1/admin/verification-records", verifyH.ListRecords)

	var h http.Handler = mux
	h = middleware.CORS(h)
	h = middleware.Logger(h)
	h = middleware.Recovery(h)
	return h
}
