package license

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/license"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

type LicenseRepository interface {
	Create(ctx context.Context, license *entity.License) error
	GetByID(ctx context.Context, licenseID uuid.UUID) (*entity.License, error)
	GetAll(ctx context.Context) ([]*entity.License, error)
	Update(ctx context.Context, license *entity.License) error
	Delete(ctx context.Context, licenseID uuid.UUID) error
}

var (
	ErrGenerateID       = errors.New("generate ID error")
	ErrInvalidLicenseID = errors.New("invalid license ID")
)

type LicenseService struct {
	repo LicenseRepository
}

func NewLicenseService(repo LicenseRepository) *LicenseService {
	return &LicenseService{repo: repo}
}

func (s *LicenseService) CreateLicense(ctx context.Context, claims *entity.Claims, license *entity.License) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrCreateLicense, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	license.ID, err = uuid.NewRandom()
	if err != nil {
		return ErrGenerateID
	}

	return s.repo.Create(ctx, license)
}

func (s *LicenseService) GetLicenseByID(ctx context.Context, licenseID uuid.UUID) (_ *entity.License, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetLicense, err)
	}()

	if licenseID == uuid.Nil {
		return nil, ErrInvalidLicenseID
	}

	return s.repo.GetByID(ctx, licenseID)
}

func (s *LicenseService) GetAllLicenses(ctx context.Context) (_ []*entity.License, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetAllLicenses, err)
	}()

	return s.repo.GetAll(ctx)
}

func (s *LicenseService) UpdateLicense(ctx context.Context, claims *entity.Claims, license *entity.License) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdateLicense, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	if license.ID == uuid.Nil {
		return ErrInvalidLicenseID
	}

	return s.repo.Update(ctx, license)
}

func (s *LicenseService) DeleteLicense(ctx context.Context, claims *entity.Claims, licenseID uuid.UUID) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrDeleteLicense, err)
	}()

	if claims == nil || claims.AccessLvl != entity.Admin {
		return commonerr.ErrForbidden
	}

	if licenseID == uuid.Nil {
		return ErrInvalidLicenseID
	}

	return s.repo.Delete(ctx, licenseID)
}
