package license_postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

type LicenseRepository struct {
	pool *pgxpool.Pool
}

func NewLicenseRepository(pool *pgxpool.Pool) *LicenseRepository {
	return &LicenseRepository{pool: pool}
}

func (r *LicenseRepository) Create(ctx context.Context, license *entity.License) error {
	query := `INSERT INTO licenses (id, title, description, url) VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, license.ID, license.Title, license.Description, license.URL)
	return err
}

func (r *LicenseRepository) GetByID(ctx context.Context, licenseID uuid.UUID) (*entity.License, error) {
	query := `SELECT id, title, description, url FROM licenses WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, licenseID)

	var l entity.License
	err := row.Scan(&l.ID, &l.Title, &l.Description, &l.URL)
	if err != nil {
		return nil, fmt.Errorf("license not found: %w", err)
	}
	return &l, nil
}

func (r *LicenseRepository) GetAll(ctx context.Context) ([]*entity.License, error) {
	query := `SELECT id, title, description, url FROM licenses`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all licenses: %w", err)
	}
	defer rows.Close()

	var licenses []*entity.License
	for rows.Next() {
		var l entity.License
		err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan license: %w", err)
		}
		licenses = append(licenses, &l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return licenses, nil
}

func (r *LicenseRepository) Update(ctx context.Context, license *entity.License) error {
	query := `UPDATE licenses SET title = $1, description = $2, url = $3 WHERE id = $4`
	_, err := r.pool.Exec(ctx, query, license.Title, license.Description, license.URL, license.ID)
	return err
}

func (r *LicenseRepository) Delete(ctx context.Context, licenseID uuid.UUID) error {
	query := `DELETE FROM licenses WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, licenseID)
	return err
}
