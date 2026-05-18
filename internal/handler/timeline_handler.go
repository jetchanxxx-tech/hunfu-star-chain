package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/huifu/star-chain/internal/store"
)

type TimelineHandler struct {
	timeline *store.TimelineStore
	family   *store.FamilyStore
}

func NewTimelineHandler(db *sql.DB) *TimelineHandler {
	return &TimelineHandler{
		timeline: store.NewTimelineStore(db),
		family:   store.NewFamilyStore(db),
	}
}

func (h *TimelineHandler) ListByMember(w http.ResponseWriter, r *http.Request) {
	memberUUID := r.PathValue("id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 50
	}

	member, err := h.family.FindMemberByUUID(memberUUID)
	if err != nil {
		Error(w, http.StatusNotFound, "member not found")
		return
	}

	events, err := h.timeline.ListByMember(member.ID, limit)
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	type item struct {
		ID        int64  `json:"id"`
		EventType string `json:"event_type"`
		EventDate string `json:"event_date"`
		EventData string `json:"event_data,omitempty"`
		Source    string `json:"source"`
	}
	var out []item
	for _, e := range events {
		out = append(out, item{
			ID:        e.ID,
			EventType: e.EventType,
			EventDate: e.EventDate.Format("2006-01-02"),
			EventData: e.EventData.String,
			Source:    e.Source,
		})
	}
	JSON(w, http.StatusOK, out)
}
