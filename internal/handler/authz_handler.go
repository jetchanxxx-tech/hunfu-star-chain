package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/huifu/star-chain/internal/service"
)

type AuthzHandler struct {
	svc *service.AuthService
}

func NewAuthzHandler(db *sql.DB) *AuthzHandler {
	return &AuthzHandler{svc: service.NewAuthService(db)}
}

func (h *AuthzHandler) RequestAuthorization(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GrantorUUID string   `json:"grantor_uuid"`
		GranteeUUID string   `json:"grantee_uuid"`
		Scopes      []string `json:"scopes"`
		Message     string   `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.GrantorUUID == "" || req.GranteeUUID == "" {
		Error(w, http.StatusBadRequest, "grantor_uuid and grantee_uuid required")
		return
	}
	if len(req.Scopes) == 0 {
		req.Scopes = []string{"report", "timeline"}
	}
	auth, err := h.svc.RequestAuthorization(req.GrantorUUID, req.GranteeUUID, req.Scopes, req.Message)
	if err != nil {
		Error(w, http.StatusConflict, err.Error())
		return
	}
	JSON(w, http.StatusCreated, map[string]any{"data": auth})
}

func (h *AuthzHandler) RespondToAuthorization(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GranteeUUID  string `json:"grantee_uuid"`
		Action       string `json:"action"` // approve / reject
		RejectReason string `json:"reject_reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	authIDStr := r.PathValue("id")
	authID, err := strconv.ParseInt(authIDStr, 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid authorization id")
		return
	}
	result, err := h.svc.RespondToAuth(authID, req.GranteeUUID, req.Action, req.RejectReason)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": result})
}

func (h *AuthzHandler) RevokeAuthorization(w http.ResponseWriter, r *http.Request) {
	var req struct {
		GrantorUUID string `json:"grantor_uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	authIDStr := r.PathValue("id")
	authID, err := strconv.ParseInt(authIDStr, 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid authorization id")
		return
	}
	if err := h.svc.RevokeAuth(authID, req.GrantorUUID); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *AuthzHandler) CheckAccess(w http.ResponseWriter, r *http.Request) {
	memberUUID := r.URL.Query().Get("member_uuid")
	targetUUID := r.URL.Query().Get("target_uuid")
	if memberUUID == "" || targetUUID == "" {
		Error(w, http.StatusBadRequest, "member_uuid and target_uuid required")
		return
	}
	allowed, err := h.svc.CheckAccess(memberUUID, targetUUID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]bool{"allowed": allowed})
}

func (h *AuthzHandler) ListAuthorizations(w http.ResponseWriter, r *http.Request) {
	memberUUID := r.URL.Query().Get("member_uuid")
	if memberUUID == "" {
		Error(w, http.StatusBadRequest, "member_uuid required")
		return
	}
	list, err := h.svc.ListMyAuthorizations(memberUUID)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": list})
}

func (h *AuthzHandler) ListAuditLogs(w http.ResponseWriter, r *http.Request) {
	offset, limit := Pagination(r)
	list, err := h.svc.ListAuditLogs(offset, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": list})
}
