package handler

import (
	"database/sql"
	"net/http"
)

type DemoHandler struct{ db *sql.DB }

func NewDemoHandler(db *sql.DB) *DemoHandler { return &DemoHandler{db: db} }

// GET /api/v1/demo/home (无需登录)
func (h *DemoHandler) Home(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		JSON(w, http.StatusOK, demoFallback())
		return
	}

	var resp struct {
		MemberUUID string           `json:"member_uuid"`
		Nickname   string           `json:"nickname"`
		Week       int              `json:"week"`
		Events     []map[string]any `json:"events"`
		Reports    []map[string]any `json:"reports"`
	}

	// 取种子数据中李娜(孕期数据最全的成员)
	row := h.db.QueryRow(`SELECT member_uuid, COALESCE(nickname,''), COALESCE(TIMESTAMPDIFF(WEEK, '2024-12-01', NOW()), 0)
		FROM family_members WHERE id = 2`)
	var uuid, nick string
	var week int
	if err := row.Scan(&uuid, &nick, &week); err != nil {
		JSON(w, http.StatusOK, demoFallback())
		return
	}
	resp.MemberUUID = uuid
	resp.Nickname = nick
	if week > 40 {
		week = 40
	}
	resp.Week = week

	// 时间轴(最近5条)
	rows, _ := h.db.Query(`SELECT event_type, DATE_FORMAT(event_date,'%Y-%m-%d'), COALESCE(event_data,'')
		FROM timeline_events WHERE member_id = 2 ORDER BY event_date DESC LIMIT 5`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var eType, eDate, eData string
			rows.Scan(&eType, &eDate, &eData)
			resp.Events = append(resp.Events, map[string]any{
				"event_type": eType, "event_date": eDate, "event_data": eData,
			})
		}
	}

	// 报告(最近3条)
	rows2, _ := h.db.Query(`SELECT report_type, COALESCE(DATE_FORMAT(report_date,'%Y-%m-%d'),''), COALESCE(summary,'')
		FROM health_reports WHERE member_id = 2 ORDER BY report_date DESC LIMIT 3`)
	if rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var rType, rDate, rSummary string
			rows2.Scan(&rType, &rDate, &rSummary)
			resp.Reports = append(resp.Reports, map[string]any{
				"report_type": rType, "report_date": rDate, "summary": rSummary,
			})
		}
	}

	JSON(w, http.StatusOK, resp)
}

func demoFallback() map[string]any {
	return map[string]any{
		"nickname": "星球居民",
		"week":     24,
		"events":   []map[string]any{},
		"reports":  []map[string]any{},
	}
}
