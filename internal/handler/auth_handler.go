package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/huifu/star-chain/internal/config"
	"github.com/huifu/star-chain/internal/service"
	"github.com/huifu/star-chain/internal/wechat"
)

type AuthHandler struct {
	wxClient *wechat.MiniprogramClient
	binding  *wechat.BindingStore
	familySvc *service.FamilyService
}

func NewAuthHandler(cfg *config.Config, db *sql.DB) *AuthHandler {
	return &AuthHandler{
		wxClient:  wechat.NewMiniprogramClient(cfg.Wechat.Miniprogram),
		binding:   wechat.NewBindingStore(db),
		familySvc: service.NewFamilyService(db),
	}
}

type wxLoginRequest struct {
	Code string `json:"code"`
}

type wxLoginResponse struct {
	Token     string `json:"token"`
	MemberID  string `json:"member_id,omitempty"`
	IsNewUser bool   `json:"is_new_user"`
}

func (h *AuthHandler) WxLogin(w http.ResponseWriter, r *http.Request) {
	var req wxLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.Code == "" {
		Error(w, http.StatusBadRequest, "code is required")
		return
	}

	session, err := h.wxClient.Code2Session(r.Context(), req.Code)
	if err != nil {
		Error(w, http.StatusUnauthorized, "wechat login failed: "+err.Error())
		return
	}

	memberID, err := h.binding.FindMemberByOpenid(session.OpenID)
	if err == sql.ErrNoRows {
		// new user, needs to register first
		JSON(w, http.StatusOK, wxLoginResponse{
			Token:     session.OpenID, // temporary: openid as token
			IsNewUser: true,
		})
		return
	}
	if err != nil {
		Error(w, http.StatusInternalServerError, "lookup failed")
		return
	}

	_ = memberID
	// existing user: generate JWT token and return
	JSON(w, http.StatusOK, wxLoginResponse{
		Token:     session.OpenID,
		IsNewUser: false,
	})
}

func (h *AuthHandler) BindWechat(w http.ResponseWriter, r *http.Request) {
	// P0 stub: bind after registration
	var req struct {
		Code     string `json:"code"`
		MemberID string `json:"member_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request")
		return
	}
	session, err := h.wxClient.Code2Session(r.Context(), req.Code)
	if err != nil {
		Error(w, http.StatusInternalServerError, "wechat error: "+err.Error())
		return
	}

	member, err := h.familySvc.FindMemberByUUID(req.MemberID)
	_ = member
	_ = err

	JSON(w, http.StatusOK, map[string]string{
		"openid":  session.OpenID,
		"unionid": session.UnionID,
	})
}
