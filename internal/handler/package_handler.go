package handler

import (
	"database/sql"
	"net/http"

	"github.com/huifu/star-chain/internal/store"
)

type PackageHandler struct {
	store *store.PackageStore
}

func NewPackageHandler(db *sql.DB) *PackageHandler {
	return &PackageHandler{store: store.NewPackageStore(db)}
}

func (h *PackageHandler) List(w http.ResponseWriter, r *http.Request) {
	packages, err := h.store.List()
	if err != nil {
		Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	type item struct {
		UUID        string  `json:"package_uuid"`
		Name        string  `json:"name"`
		Description string  `json:"description,omitempty"`
		Level       string  `json:"level"`
		Price       float64 `json:"price"`
		CoverImage  string  `json:"cover_image,omitempty"`
		Benefits    string  `json:"benefits,omitempty"`
		Status      string  `json:"status"`
	}
	var out []item
	for _, p := range packages {
		out = append(out, item{
			UUID:        p.UUID,
			Name:        p.Name,
			Description: p.Description.String,
			Level:       p.Level,
			Price:       p.Price,
			CoverImage:  p.CoverImage.String,
			Benefits:    p.Benefits.String,
			Status:      p.Status,
		})
	}
	JSON(w, http.StatusOK, out)
}
