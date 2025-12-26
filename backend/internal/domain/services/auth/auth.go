package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/auth"
	"github.com/hahaclassic/orpheon/backend/pkg/errwrap"
)

var (
	ErrWrongOldPassword   = errors.New("wrong old password")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type TokenService interface {
	GenerateAccessToken(claims *entity.Claims) (string, error)
	ParseAccessToken(tokenStr string) (*entity.Claims, error)
	GenerateRefreshToken() (string, error)
}

type PasswordHasher interface {
	GenerateFromPassword(password string) (string, error)
	CompareHashAndPassword(hashed string, password string) error
}

type RefreshTokenRepository interface {
	Set(ctx context.Context, token string, claims *entity.Claims) error
	Get(ctx context.Context, token string) (*entity.Claims, error)
	Delete(ctx context.Context, token string) error
}

type AuthRepository interface {
	SaveCredentials(ctx context.Context, userID uuid.UUID, credentials *entity.UserCredentials) error
	GetPasswordByLogin(ctx context.Context, login string) (string, error)
	GetPasswordByID(ctx context.Context, userID uuid.UUID) (string, error)
	GetClaimsByLogin(ctx context.Context, login string) (*entity.Claims, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
}

type UserCreatorService interface {
	CreateUser(ctx context.Context, info *entity.UserInfo) (uuid.UUID, error)
}

type AuthService struct {
	authRepo     AuthRepository
	refreshRepo  RefreshTokenRepository
	userCreator  UserCreatorService
	hasher       PasswordHasher
	tokenService TokenService
}

func NewAuthService(authRepo AuthRepository, refreshRepo RefreshTokenRepository,
	userCreator UserCreatorService, hasher PasswordHasher, token TokenService) *AuthService {
	return &AuthService{
		authRepo:     authRepo,
		refreshRepo:  refreshRepo,
		userCreator:  userCreator,
		hasher:       hasher,
		tokenService: token,
	}
}

func (a *AuthService) RegisterUser(ctx context.Context, credentials *entity.UserCredentials) (_ *entity.AuthTokens, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrRegisterUser, err)
	}()

	hashedPassword, err := a.hasher.GenerateFromPassword(credentials.Password)
	if err != nil {
		return nil, err
	}

	userID, err := a.userCreator.CreateUser(ctx, &entity.UserInfo{Name: credentials.Login})
	if err != nil {
		return nil, err
	}

	hashedCreds := &entity.UserCredentials{
		Login:    credentials.Login,
		Password: string(hashedPassword),
	}

	err = a.authRepo.SaveCredentials(ctx, userID, hashedCreds)
	if err != nil {
		return nil, err
	}

	return a.Login(ctx, credentials)
}

func (a *AuthService) Login(ctx context.Context, credentials *entity.UserCredentials) (_ *entity.AuthTokens, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrLogin, err)
	}()

	hashedPassword, err := a.authRepo.GetPasswordByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, err
	}

	if a.hasher.CompareHashAndPassword(hashedPassword, credentials.Password) != nil {
		return nil, ErrInvalidCredentials
	}

	claims, err := a.authRepo.GetClaimsByLogin(ctx, credentials.Login)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.tokenService.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = a.refreshRepo.Set(ctx, refreshToken, claims)
	if err != nil {
		return nil, err
	}

	return &entity.AuthTokens{
		Access:  accessToken,
		Refresh: refreshToken}, nil
}

func (a *AuthService) Logout(ctx context.Context, refreshToken string) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrLogout, err)
	}()

	return a.refreshRepo.Delete(ctx, refreshToken)
}

func (a *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (_ *entity.AuthTokens, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrRefreshTokens, err)
	}()

	claims, err := a.refreshRepo.Get(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.tokenService.GenerateAccessToken(claims)
	if err != nil {
		return nil, err
	}

	newRefresh, err := a.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	err = a.refreshRepo.Delete(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	err = a.refreshRepo.Set(ctx, newRefresh, claims)
	if err != nil {
		return nil, err
	}

	return &entity.AuthTokens{Access: accessToken, Refresh: newRefresh}, nil
}

func (a *AuthService) UpdatePassword(ctx context.Context, userID uuid.UUID, passwords *entity.UserPasswords) (err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrUpdatePassword, err)
	}()

	hashedPassword, err := a.authRepo.GetPasswordByID(ctx, userID)
	if err != nil {
		return err
	}

	if a.hasher.CompareHashAndPassword(hashedPassword, passwords.Old) != nil {
		return ErrInvalidCredentials
	}

	newHashed, err := a.hasher.GenerateFromPassword(passwords.New)
	if err != nil {
		return err
	}

	return a.authRepo.UpdatePassword(ctx, userID, string(newHashed))
}

func (a *AuthService) GetClaims(ctx context.Context, accessToken string) (_ *entity.Claims, err error) {
	defer func() {
		err = errwrap.WrapIfErr(usecase.ErrGetClaims, err)
	}()

	return a.tokenService.ParseAccessToken(accessToken)
}
