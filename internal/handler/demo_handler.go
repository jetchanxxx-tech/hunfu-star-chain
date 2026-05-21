package handler

import (
	"database/sql"
	"net/http"
)

type DemoHandler struct{ db *sql.DB }

func NewDemoHandler(db *sql.DB) *DemoHandler { return &DemoHandler{db: db} }

var eventCN = map[string]string{
	"first_prenatal": "首次产检", "nt": "NT检查", "early_tang": "早唐筛查",
	"ogtt": "糖耐量", "quad_d": "四维彩超", "delivery": "分娩", "42day": "产后42天复查",
	"vaccine_2m": "2月龄疫苗", "vaccine_3m": "3月龄疫苗", "vaccine": "疫苗接种",
}

var reportCN = map[string]string{
	"lab": "检验报告", "imaging": "检查报告", "discharge": "出院小结",
}

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
		Stats      map[string]int   `json:"stats"`
		Events     []map[string]any `json:"events"`
		Reports    []map[string]any `json:"reports"`
		Packages   []map[string]any `json:"packages"`
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
	if week > 40 { week = 40 }
	resp.Week = week

	// 统计
	var memberCount, familyCount, pkgCount int
	h.db.QueryRow(`SELECT COUNT(*) FROM family_members`).Scan(&memberCount)
	h.db.QueryRow(`SELECT COUNT(*) FROM families`).Scan(&familyCount)
	h.db.QueryRow(`SELECT COUNT(*) FROM service_packages WHERE status='online'`).Scan(&pkgCount)
	resp.Stats = map[string]int{"members": memberCount, "families": familyCount, "packages": pkgCount}

	// 时间轴(最近5条) - 中文标签
	rows, _ := h.db.Query(`SELECT event_type, DATE_FORMAT(event_date,'%Y-%m-%d'), COALESCE(event_data,'')
		FROM timeline_events WHERE member_id = 2 ORDER BY event_date DESC LIMIT 5`)
	if rows != nil {
		defer rows.Close()
		for rows.Next() {
			var eType, eDate, eData string
			rows.Scan(&eType, &eDate, &eData)
			label := eventCN[eType]
			if label == "" { label = eType }
			resp.Events = append(resp.Events, map[string]any{
				"event_type": label, "event_date": eDate, "event_data": eData,
			})
		}
	}

	// 报告(最近3条) - 中文标签
	rows2, _ := h.db.Query(`SELECT report_type, COALESCE(DATE_FORMAT(report_date,'%Y-%m-%d'),''), COALESCE(summary,'')
		FROM health_reports WHERE member_id = 2 ORDER BY report_date DESC LIMIT 3`)
	if rows2 != nil {
		defer rows2.Close()
		for rows2.Next() {
			var rType, rDate, rSummary string
			rows2.Scan(&rType, &rDate, &rSummary)
			label := reportCN[rType]
			if label == "" { label = rType }
			resp.Reports = append(resp.Reports, map[string]any{
				"report_type": label, "report_date": rDate, "summary": rSummary,
			})
		}
	}

	// 服务包(前3个 online)
	rows3, _ := h.db.Query(`SELECT name, level, price, COALESCE(description,''), COALESCE(benefits,'')
		FROM service_packages WHERE status = 'online' ORDER BY level, price LIMIT 3`)
	if rows3 != nil {
		defer rows3.Close()
		for rows3.Next() {
			var name, level, desc, benefits string
			var price float64
			rows3.Scan(&name, &level, &price, &desc, &benefits)
			resp.Packages = append(resp.Packages, map[string]any{
				"name": name, "level": level, "price": price, "description": desc, "benefits": benefits,
			})
		}
	}

	JSON(w, http.StatusOK, resp)
}

func demoFallback() map[string]any {
	return map[string]any{
		"nickname": "星球居民", "week": 24,
		"stats":    map[string]int{"members": 0, "families": 0, "packages": 0},
		"events":   []map[string]any{}, "reports": []map[string]any{}, "packages": []map[string]any{},
	}
}
