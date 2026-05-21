package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/huifu/star-chain/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
	store *store.AdminStore
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{store: store.NewAdminStore(db)}
}

// POST /api/v1/admin/login
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, http.StatusBadRequest, "invalid request")
		return
	}
	user, err := h.store.FindAdminByUsername(req.Username)
	if err != nil || user.Status != "active" {
		Error(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		Error(w, http.StatusUnauthorized, "用户名或密码错误")
		return
	}
	h.store.UpdateLastLogin(user.ID)
	JSON(w, http.StatusOK, map[string]any{
		"token":    "admin-session-" + user.Username,
		"username": user.Username,
		"role":     user.Role,
		"real_name": user.RealName,
	})
}

// GET /api/v1/admin/dashboard
func (h *AdminHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetDashboardStats()
	if err != nil {
		Error(w, http.StatusInternalServerError, "query failed")
		return
	}
	JSON(w, http.StatusOK, stats)
}

// GET /api/v1/admin/members
func (h *AdminHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	members, err := h.store.ListAllMembers(search)
	if err != nil {
		Error(w, http.StatusInternalServerError, "query failed")
		return
	}
	if members == nil {
		members = []store.MemberRow{}
	}
	JSON(w, http.StatusOK, members)
}

// GET /api/v1/admin/packages
func (h *AdminHandler) ListPackages(w http.ResponseWriter, r *http.Request) {
	pkgs, err := h.store.ListPackages()
	if err != nil {
		Error(w, http.StatusInternalServerError, "query failed")
		return
	}
	if pkgs == nil {
		pkgs = []store.PackageRow{}
	}
	JSON(w, http.StatusOK, pkgs)
}

// POST /api/v1/admin/packages
func (h *AdminHandler) CreatePackage(w http.ResponseWriter, r *http.Request) {
	var pkg store.PackageRow
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		Error(w, http.StatusBadRequest, "invalid request")
		return
	}
	if err := h.store.CreatePackage(&pkg); err != nil {
		Error(w, http.StatusInternalServerError, "create failed: "+err.Error())
		return
	}
	JSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

// PUT /api/v1/admin/packages/{id}
func (h *AdminHandler) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var pkg store.PackageRow
	if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
		Error(w, http.StatusBadRequest, "invalid request")
		return
	}
	// Parse id from path
	pid := 0
	for _, c := range id {
		pid = pid*10 + int(c-'0')
	}
	pkg.ID = int64(pid)
	if err := h.store.UpdatePackage(&pkg); err != nil {
		Error(w, http.StatusInternalServerError, "update failed: "+err.Error())
		return
	}
	JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// GET /api/v1/admin/users
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.ListAdminUsers()
	if err != nil {
		Error(w, http.StatusInternalServerError, "query failed")
		return
	}
	if users == nil {
		users = []store.AdminUser{}
	}
	JSON(w, http.StatusOK, users)
}
