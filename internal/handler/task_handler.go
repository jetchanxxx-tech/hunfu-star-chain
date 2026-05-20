package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/huifu/star-chain/internal/service"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{svc: service.NewTaskService(db)}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MemberID     int64  `json:"member_id"`
		TriggerType  string `json:"trigger_type"`
		TriggerValue string `json:"trigger_value"`
		Title        string `json:"title"`
		DueDate      string `json:"due_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Title == "" {
		Error(w, http.StatusBadRequest, "title required")
		return
	}
	dueDate, _ := time.Parse("2006-01-02", req.DueDate)
	task, err := h.svc.CreateTask(req.MemberID, req.TriggerType, req.TriggerValue, req.Title, dueDate)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusCreated, map[string]any{"data": task})
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filters := map[string]interface{}{}
	if v := q.Get("status"); v != "" {
		filters["status"] = v
	}
	if v := q.Get("trigger_type"); v != "" {
		filters["trigger_type"] = v
	}
	if v := q.Get("assigned_to"); v != "" {
		id, _ := strconv.ParseInt(v, 10, 64)
		filters["assigned_to"] = id
	}
	if v := q.Get("member_id"); v != "" {
		id, _ := strconv.ParseInt(v, 10, 64)
		filters["member_id"] = id
	}
	offset, limit := Pagination(r)
	list, err := h.svc.ListTasks(filters, offset, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": list})
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	task, err := h.svc.GetTask(id)
	if err != nil {
		Error(w, http.StatusNotFound, "task not found")
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": task})
}

func (h *TaskHandler) Assign(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req struct {
		StewardID int64 `json:"steward_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	if err := h.svc.AssignTask(id, req.StewardID); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *TaskHandler) Complete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req struct {
		Notes string `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if err := h.svc.CompleteTask(id, req.Notes); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *TaskHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	var req struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Reason == "" {
		Error(w, http.StatusBadRequest, "reason required")
		return
	}
	if err := h.svc.CancelTask(id, req.Reason); err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *TaskHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.svc.Stats()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]any{"data": stats})
}
