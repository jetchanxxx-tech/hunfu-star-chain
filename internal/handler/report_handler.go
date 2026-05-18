package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/huifu/star-chain/internal/service"
)

type ReportHandler struct {
	svc *service.ReportService
}

func NewReportHandler(db *sql.DB) *ReportHandler {
	return &ReportHandler{svc: service.NewReportService(db)}
}

func (h *ReportHandler) ListByMember(w http.ResponseWriter, r *http.Request) {
	memberUUID := r.PathValue("id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

	reports, err := h.svc.ListByMember(memberUUID, limit, offset)
	if err != nil {
		Error(w, http.StatusNotFound, err.Error())
		return
	}
	JSON(w, http.StatusOK, reports)
}
