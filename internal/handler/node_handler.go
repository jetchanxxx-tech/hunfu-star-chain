package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/huifu/star-chain/internal/model"
	"github.com/huifu/star-chain/internal/service"
)

type NodeHandler struct {
	svc *service.NodeService
}

func NewNodeHandler(db *sql.DB) *NodeHandler {
	return &NodeHandler{svc: service.NewNodeService(db)}
}

func (h *NodeHandler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := h.svc.ListTemplates()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": templates})
}

func (h *NodeHandler) UpsertTemplate(w http.ResponseWriter, r *http.Request) {
	var t model.TimelineNodeTemplate
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if t.NodeCode == "" || t.NodeName == "" {
		Error(w, http.StatusBadRequest, "node_code and node_name are required")
		return
	}
	if err := h.svc.UpsertTemplate(&t); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *NodeHandler) UpdateTemplateStatus(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.svc.UpdateTemplateStatus(code, req.Status); err != nil {
		Error(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *NodeHandler) ListOverrides(w http.ResponseWriter, r *http.Request) {
	hospitalCode := r.PathValue("hospitalCode")
	overrides, err := h.svc.ListOverrides(hospitalCode)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": overrides})
}

func (h *NodeHandler) UpsertOverride(w http.ResponseWriter, r *http.Request) {
	var o model.HospitalNodeOverride
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if o.HospitalCode == "" || o.NodeCode == "" {
		Error(w, http.StatusBadRequest, "hospital_code and node_code are required")
		return
	}
	if err := h.svc.UpsertOverride(&o); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *NodeHandler) DeleteOverride(w http.ResponseWriter, r *http.Request) {
	hospitalCode := r.PathValue("hospitalCode")
	nodeCode := r.URL.Query().Get("node_code")
	if nodeCode == "" {
		Error(w, http.StatusBadRequest, "node_code query param required")
		return
	}
	if err := h.svc.DeleteOverride(hospitalCode, nodeCode); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
