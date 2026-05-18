package store

import (
	"database/sql"

	"github.com/huifu/star-chain/internal/model"
)

type PackageStore struct{ db *sql.DB }

func NewPackageStore(db *sql.DB) *PackageStore { return &PackageStore{db: db} }

func (s *PackageStore) List() ([]model.ServicePackage, error) {
	rows, err := s.db.Query(
		`SELECT id, package_uuid, name, description, level, price, cover_image, benefits, status, created_at, updated_at
		 FROM service_packages WHERE status = 'online' ORDER BY level, price`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var packages []model.ServicePackage
	for rows.Next() {
		var p model.ServicePackage
		if err := rows.Scan(&p.ID, &p.UUID, &p.Name, &p.Description, &p.Level, &p.Price,
			&p.CoverImage, &p.Benefits, &p.Status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		packages = append(packages, p)
	}
	return packages, rows.Err()
}

func (s *PackageStore) FindByUUID(uuid string) (*model.ServicePackage, error) {
	var p model.ServicePackage
	err := s.db.QueryRow(
		`SELECT id, package_uuid, name, description, level, price, cover_image, benefits, status, created_at, updated_at
		 FROM service_packages WHERE package_uuid = ?`, uuid,
	).Scan(&p.ID, &p.UUID, &p.Name, &p.Description, &p.Level, &p.Price,
		&p.CoverImage, &p.Benefits, &p.Status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
