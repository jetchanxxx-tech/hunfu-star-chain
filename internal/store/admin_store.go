package store

import (
	"database/sql"
	"time"
)

type AdminStore struct{ db *sql.DB }

func NewAdminStore(db *sql.DB) *AdminStore { return &AdminStore{db: db} }

// --- Login ---
type AdminUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
	RealName string `json:"real_name"`
	Status   string `json:"status"`
}

func (s *AdminStore) FindAdminByUsername(username string) (*AdminUser, error) {
	var u AdminUser
	err := s.db.QueryRow(
		`SELECT id, username, password, role, real_name, status FROM admin_users WHERE username = ?`, username,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.RealName, &u.Status)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *AdminStore) UpdateLastLogin(id int64) {
	s.db.Exec(`UPDATE admin_users SET last_login = NOW() WHERE id = ?`, id)
}

func (s *AdminStore) ListAdminUsers() ([]AdminUser, error) {
	rows, err := s.db.Query(`SELECT id, username, role, real_name, status FROM admin_users ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []AdminUser
	for rows.Next() {
		var u AdminUser
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.RealName, &u.Status); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// --- Dashboard ---
type DashboardStats struct {
	TotalMembers    int `json:"total_members"`
	NewMembersMonth int `json:"new_members_month"`
	ActivePackages  int `json:"active_packages"`
	TaskCompleteRate float64 `json:"task_complete_rate"`
	PendingTasks    int `json:"pending_tasks"`
	TotalFamilies   int `json:"total_families"`
	ActiveMembers   int `json:"active_members"`
	InactiveMembers int `json:"inactive_members"`
}

func (s *AdminStore) GetDashboardStats() (*DashboardStats, error) {
	var ds DashboardStats
	// Total members
	s.db.QueryRow(`SELECT COUNT(*) FROM family_members`).Scan(&ds.TotalMembers)
	// Active members
	s.db.QueryRow(`SELECT COUNT(*) FROM family_members WHERE status = 'active'`).Scan(&ds.ActiveMembers)
	// Inactive
	s.db.QueryRow(`SELECT COUNT(*) FROM family_members WHERE status = 'inactive'`).Scan(&ds.InactiveMembers)
	// Families
	s.db.QueryRow(`SELECT COUNT(*) FROM families`).Scan(&ds.TotalFamilies)
	// New members this month
	s.db.QueryRow(`SELECT COUNT(*) FROM family_members WHERE created_at >= DATE_SUB(NOW(), INTERVAL 1 MONTH)`).Scan(&ds.NewMembersMonth)
	// Active packages
	s.db.QueryRow(`SELECT COUNT(*) FROM member_packages WHERE status = 'active'`).Scan(&ds.ActivePackages)
	// Pending tasks
	s.db.QueryRow(`SELECT COUNT(*) FROM followup_tasks WHERE status IN ('pending','in_progress')`).Scan(&ds.PendingTasks)
	// Task complete rate
	var completed, total int
	s.db.QueryRow(`SELECT COUNT(*) FROM followup_tasks WHERE status = 'completed'`).Scan(&completed)
	s.db.QueryRow(`SELECT COUNT(*) FROM followup_tasks`).Scan(&total)
	if total > 0 {
		ds.TaskCompleteRate = float64(completed) / float64(total) * 100
	}
	return &ds, nil
}

// --- Members ---
type MemberRow struct {
	ID         int64     `json:"id"`
	MemberUUID string    `json:"member_uuid"`
	Nickname   string    `json:"nickname"`
	Relation   string    `json:"relation"`
	Gender     int       `json:"gender"`
	BirthDate  string    `json:"birth_date"`
	FamilyName string    `json:"family_name"`
	PackageName string   `json:"package_name"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

func (s *AdminStore) ListAllMembers(search string) ([]MemberRow, error) {
	query := `SELECT m.id, m.member_uuid, COALESCE(m.nickname,''), m.relation,
		COALESCE(m.gender,0), COALESCE(DATE_FORMAT(m.birth_date,'%Y-%m-%d'),''),
		COALESCE(f.name,''), COALESCE(sp.name,'无'), m.status, m.created_at
		FROM family_members m
		LEFT JOIN families f ON f.id = m.family_id
		LEFT JOIN member_packages mp ON mp.member_id = m.id AND mp.status = 'active'
		LEFT JOIN service_packages sp ON sp.id = mp.package_id`
	args := []any{}
	if search != "" {
		query += ` WHERE m.nickname LIKE ? OR f.name LIKE ?`
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	query += ` ORDER BY m.created_at DESC LIMIT 200`
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var members []MemberRow
	for rows.Next() {
		var m MemberRow
		if err := rows.Scan(&m.ID, &m.MemberUUID, &m.Nickname, &m.Relation,
			&m.Gender, &m.BirthDate, &m.FamilyName, &m.PackageName, &m.Status, &m.CreatedAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

// --- Service Packages ---
type PackageRow struct {
	ID          int64     `json:"id"`
	PackageUUID string    `json:"package_uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Level       string    `json:"level"`
	Price       float64   `json:"price"`
	Benefits    string    `json:"benefits"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *AdminStore) ListPackages() ([]PackageRow, error) {
	rows, err := s.db.Query(`SELECT id, package_uuid, name, COALESCE(description,''), level, price,
		COALESCE(benefits,''), status, created_at FROM service_packages ORDER BY level, price`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pkgs []PackageRow
	for rows.Next() {
		var p PackageRow
		if err := rows.Scan(&p.ID, &p.PackageUUID, &p.Name, &p.Description,
			&p.Level, &p.Price, &p.Benefits, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		pkgs = append(pkgs, p)
	}
	return pkgs, rows.Err()
}

func (s *AdminStore) CreatePackage(p *PackageRow) error {
	_, err := s.db.Exec(`INSERT INTO service_packages (package_uuid, name, description, level, price, benefits, status)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, p.PackageUUID, p.Name, p.Description, p.Level, p.Price, p.Benefits, p.Status)
	return err
}

func (s *AdminStore) UpdatePackage(p *PackageRow) error {
	_, err := s.db.Exec(`UPDATE service_packages SET name=?, description=?, level=?, price=?, benefits=?, status=? WHERE id=?`,
		p.Name, p.Description, p.Level, p.Price, p.Benefits, p.Status, p.ID)
	return err
}
