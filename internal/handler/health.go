package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Health struct{ db *sql.DB }

func NewHealth(db *sql.DB) *Health { return &Health{db: db} }

func (h *Health) Check(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"status": "ok"}
	if h.db != nil {
		if err := h.db.Ping(); err != nil {
			resp["mysql"] = "error: " + err.Error()
		} else {
			resp["mysql"] = "ok"
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
