package license

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
)

var (
	ErrGetLicense     = errors.New("failed to get license")
	ErrGetAllLicenses = errors.New("failed to get all licenses")
	ErrCreateLicense  = errors.New("failed to create license")
	ErrUpdateLicense  = errors.New("failed to update license")
	ErrDeleteLicense  = errors.New("failed to delete license")
)

type LicenseService interface {
	GetLicenseByID(ctx context.Context, licenseID uuid.UUID) (_ *entity.License, err error)
	GetAllLicenses(ctx context.Context) (_ []*entity.License, err error)

	CreateLicense(ctx context.Context, claims *entity.Claims, license *entity.License) (err error)
	UpdateLicense(ctx context.Context, claims *entity.Claims, license *entity.License) (err error)
	DeleteLicense(ctx context.Context, claims *entity.Claims, licenseID uuid.UUID) (err error)
}
