package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/huifu/star-chain/internal/model"
	"github.com/huifu/star-chain/internal/service"
)

type FamilyHandler struct {
	svc *service.FamilyService
}

func NewFamilyHandler(db *sql.DB) *FamilyHandler {
	return &FamilyHandler{svc: service.NewFamilyService(db)}
}

func (h *FamilyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req model.CreateFamilyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PhoneHash == "" {
		Error(w, http.StatusBadRequest, "phone is required")
		return
	}
	resp, err := h.svc.Create(req)
	if err != nil {
		Error(w, http.StatusConflict, err.Error())
		return
	}
	JSON(w, http.StatusCreated, resp)
}

func (h *FamilyHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	resp, err := h.svc.GetFamily(id)
	if err != nil {
		Error(w, http.StatusNotFound, "family not found")
		return
	}
	JSON(w, http.StatusOK, resp)
}

func (h *FamilyHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	familyID := r.PathValue("id")
	var req model.AddMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.PhoneHash == "" {
		Error(w, http.StatusBadRequest, "phone is required")
		return
	}
	resp, err := h.svc.AddMember(familyID, req)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, resp)
}
