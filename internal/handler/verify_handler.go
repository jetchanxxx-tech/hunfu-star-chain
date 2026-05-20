package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/huifu/star-chain/internal/service"
)

type VerifyHandler struct {
	svc *service.VerifyService
}

func NewVerifyHandler(db *sql.DB, secret string) *VerifyHandler {
	return &VerifyHandler{svc: service.NewVerifyService(db, secret)}
}

func (h *VerifyHandler) GenerateQR(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MemberUUID    string `json:"member_uuid"`
		EntitlementID int64  `json:"entitlement_id"`
		BenefitType   string `json:"benefit_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.MemberUUID == "" || req.EntitlementID == 0 {
		Error(w, http.StatusBadRequest, "member_uuid and entitlement_id required")
		return
	}
	payload, err := h.svc.GenerateQRPayload(req.MemberUUID, req.EntitlementID, req.BenefitType)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": payload})
}

func (h *VerifyHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Payload   string `json:"payload"` // JSON-encoded QRPayload
		StewardID int64  `json:"steward_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	payload, err := service.DecodePayload(req.Payload)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid payload format")
		return
	}
	record, err := h.svc.Verify(payload, req.StewardID)
	if err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": record})
}

func (h *VerifyHandler) ListRecords(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("member_id")
	var memberID int64
	if idStr != "" {
		if v, err := strconv.ParseInt(idStr, 10, 64); err == nil {
			memberID = v
		}
	}
	offset, limit := Pagination(r)
	list, err := h.svc.ListRecords(memberID, offset, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": list})
}
