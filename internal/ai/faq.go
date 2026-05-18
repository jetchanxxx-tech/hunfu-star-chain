package ai

import (
	"database/sql"
	"strings"
)

type FAQEntry struct {
	ID       int64
	Category string
	Question string
	Answer   string
	Keywords string
	Priority int
}

type FAQMatcher struct{ db *sql.DB }

func NewFAQMatcher(db *sql.DB) *FAQMatcher { return &FAQMatcher{db: db} }

func (m *FAQMatcher) Match(query string) (*FAQEntry, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	rows, err := m.db.Query(
		`SELECT id, category, question, answer, keywords, priority
		 FROM faq_entries WHERE status = 'published'
		 ORDER BY priority DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var best *FAQEntry
	var bestScore int
	for rows.Next() {
		var e FAQEntry
		if err := rows.Scan(&e.ID, &e.Category, &e.Question, &e.Answer, &e.Keywords, &e.Priority); err != nil {
			continue
		}
		score := matchScore(query, e.Question, e.Keywords)
		if score > bestScore {
			bestScore = score
			entry := e
			best = &entry
		}
	}
	if bestScore < 2 {
		return nil, nil // no good match
	}
	return best, rows.Err()
}

func matchScore(query, question, keywords string) int {
	score := 0
	ql := strings.ToLower(query)
	qlWords := strings.Fields(ql)

	// keyword match
	for _, kw := range strings.Split(keywords, ",") {
		kw = strings.TrimSpace(strings.ToLower(kw))
		if kw != "" && strings.Contains(ql, kw) {
			score += 3
		}
	}

	// question word overlap
	qWords := strings.Fields(strings.ToLower(question))
	for _, qw := range qlWords {
		for _, fw := range qWords {
			if qw == fw {
				score++
			}
		}
	}
	return score
}

func (m *FAQMatcher) ListCategories() ([]string, error) {
	rows, err := m.db.Query(
		`SELECT DISTINCT category FROM faq_entries WHERE status = 'published' ORDER BY category`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cats []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			continue
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}
