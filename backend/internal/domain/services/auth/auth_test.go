package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestAuthServiceSuite(t *testing.T) {
	suite.Run(t, &AuthServiceSuite{})
}

type AuthObjectMother struct{}

func (AuthObjectMother) DefaultCredentials() *entity.UserCredentials {
	return &entity.UserCredentials{Login: "testuser", Password: "securepass"}
}

func (a AuthObjectMother) DefaultClaims() *entity.Claims {
	return &entity.Claims{UserID: a.DefaultUserID(), AccessLvl: entity.User}
}

func (AuthObjectMother) DefaultTokens() *entity.AuthTokens {
	return &entity.AuthTokens{Access: "access.token", Refresh: "refresh.token"}
}

func (AuthObjectMother) DefaultHashedPassword() string {
	return "hashed_password"
}

func (AuthObjectMother) DefaultUserID() uuid.UUID {
	return uuid.MustParse("8580f5cb-6b75-4046-8d31-6e0fe79f12c6")
}

func (AuthObjectMother) DefaultPasswords() *entity.UserPasswords {
	return &entity.UserPasswords{
		Old: "old",
		New: "new",
	}
}

type AuthServiceSuite struct {
	suite.Suite

	ctx          context.Context
	service      *AuthService
	authRepo     *mocks.AuthRepository
	refreshRepo  *mocks.RefreshTokenRepository
	userCreator  *mocks.UserCreatorService
	hasher       *mocks.PasswordHasher
	tokenService *mocks.TokenService

	objMother *AuthObjectMother
}

func (s *AuthServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.authRepo = mocks.NewAuthRepository(s.T())
	s.refreshRepo = mocks.NewRefreshTokenRepository(s.T())
	s.userCreator = mocks.NewUserCreatorService(s.T())
	s.hasher = mocks.NewPasswordHasher(s.T())
	s.tokenService = mocks.NewTokenService(s.T())

	s.service = NewAuthService(
		s.authRepo,
		s.refreshRepo,
		s.userCreator,
		s.hasher,
		s.tokenService,
	)

	s.objMother = &AuthObjectMother{}
}

// --- RegisterUser ---

func (s *AuthServiceSuite) TestRegisterUser_Success() {
	creds := s.objMother.DefaultCredentials()
	userID := s.objMother.DefaultUserID()
	hashedPassword := s.objMother.DefaultHashedPassword()
	claims := s.objMother.DefaultClaims()
	tokens := s.objMother.DefaultTokens()

	s.hasher.On("GenerateFromPassword", creds.Password).Return(hashedPassword, nil)
	s.userCreator.On("CreateUser", s.ctx, mock.MatchedBy(func(info *entity.UserInfo) bool {
		return info.Name == creds.Login
	})).Return(userID, nil)

	s.authRepo.On("SaveCredentials", s.ctx, userID, &entity.UserCredentials{
		Login:    creds.Login,
		Password: hashedPassword,
	}).Return(nil)

	s.authRepo.On("GetPasswordByLogin", s.ctx, creds.Login).Return(hashedPassword, nil)
	s.hasher.On("CompareHashAndPassword", hashedPassword, creds.Password).Return(nil)
	s.authRepo.On("GetClaimsByLogin", s.ctx, creds.Login).Return(claims, nil)
	s.tokenService.On("GenerateAccessToken", claims).Return(tokens.Access, nil)
	s.tokenService.On("GenerateRefreshToken").Return(tokens.Refresh, nil)
	s.refreshRepo.On("Set", s.ctx, tokens.Refresh, claims).Return(nil)

	resulTtokens, err := s.service.RegisterUser(s.ctx, creds)

	s.Require().NoError(err)
	s.Equal(tokens.Access, resulTtokens.Access)
	s.Equal(tokens.Refresh, resulTtokens.Refresh)

	s.hasher.AssertExpectations(s.T())
	s.userCreator.AssertExpectations(s.T())
	s.authRepo.AssertExpectations(s.T())
	s.tokenService.AssertExpectations(s.T())
	s.refreshRepo.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestRegisterUser_HasherError() {
	creds := s.objMother.DefaultCredentials()
	s.hasher.On("GenerateFromPassword", creds.Password).Return("", errors.New("hash error"))

	tokens, err := s.service.RegisterUser(s.ctx, creds)

	s.Error(err)
	s.Nil(tokens)
	s.hasher.AssertExpectations(s.T())
}

// --- Login ---

func (s *AuthServiceSuite) TestLogin_Success() {
	creds := s.objMother.DefaultCredentials()
	hashed := s.objMother.DefaultHashedPassword()
	claims := s.objMother.DefaultClaims()
	tokens := s.objMother.DefaultTokens()

	s.authRepo.On("GetPasswordByLogin", s.ctx, creds.Login).Return(hashed, nil)
	s.hasher.On("CompareHashAndPassword", hashed, creds.Password).Return(nil)
	s.authRepo.On("GetClaimsByLogin", s.ctx, creds.Login).Return(claims, nil)
	s.tokenService.On("GenerateAccessToken", claims).Return(tokens.Access, nil)
	s.tokenService.On("GenerateRefreshToken").Return(tokens.Refresh, nil)
	s.refreshRepo.On("Set", s.ctx, tokens.Refresh, claims).Return(nil)

	resultTokens, err := s.service.Login(s.ctx, creds)

	s.NoError(err)
	s.Equal(tokens, resultTokens)

	s.authRepo.AssertExpectations(s.T())
	s.hasher.AssertExpectations(s.T())
	s.tokenService.AssertExpectations(s.T())
	s.refreshRepo.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestLogin_GetPasswordError() {
	creds := s.objMother.DefaultCredentials()
	s.authRepo.On("GetPasswordByLogin", s.ctx, creds.Login).Return("", errors.New("db error"))

	tokens, err := s.service.Login(s.ctx, creds)

	s.Error(err)
	s.Nil(tokens)
	s.authRepo.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestLogin_InvalidPassword() {
	creds := s.objMother.DefaultCredentials()
	hashed := s.objMother.DefaultHashedPassword()

	s.authRepo.On("GetPasswordByLogin", s.ctx, creds.Login).Return(hashed, nil)
	s.hasher.On("CompareHashAndPassword", hashed, creds.Password).Return(errors.New("mismatch"))

	tokens, err := s.service.Login(s.ctx, creds)

	s.Error(err)
	s.Nil(tokens)
	s.ErrorIs(err, ErrInvalidCredentials)
	s.authRepo.AssertExpectations(s.T())
	s.hasher.AssertExpectations(s.T())
}

// --- Logout ---

func (s *AuthServiceSuite) TestLogout_Success() {
	refreshToken := "refresh"
	s.refreshRepo.On("Delete", s.ctx, refreshToken).Return(nil)

	err := s.service.Logout(s.ctx, refreshToken)
	s.NoError(err)
	s.refreshRepo.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestLogout_DeleteError() {
	refreshToken := "refresh"
	s.refreshRepo.On("Delete", s.ctx, refreshToken).Return(errors.New("delete error"))

	err := s.service.Logout(s.ctx, refreshToken)
	s.Error(err)
	s.refreshRepo.AssertExpectations(s.T())
}

// --- RefreshTokens ---

func (s *AuthServiceSuite) TestRefreshTokens_Success() {
	claims := s.objMother.DefaultClaims()
	oldRefresh := "old_refresh"
	newRefresh := "new_refresh"
	access := "access"

	s.refreshRepo.On("Get", s.ctx, oldRefresh).Return(claims, nil)
	s.tokenService.On("GenerateAccessToken", claims).Return(access, nil)
	s.tokenService.On("GenerateRefreshToken").Return(newRefresh, nil)
	s.refreshRepo.On("Delete", s.ctx, oldRefresh).Return(nil)
	s.refreshRepo.On("Set", s.ctx, newRefresh, claims).Return(nil)

	tokens, err := s.service.RefreshTokens(s.ctx, oldRefresh)

	s.NoError(err)
	s.Equal(&entity.AuthTokens{Access: access, Refresh: newRefresh}, tokens)
	s.refreshRepo.AssertExpectations(s.T())
	s.tokenService.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestRefreshTokens_GetError() {
	oldRefresh := "old"
	s.refreshRepo.On("Get", s.ctx, oldRefresh).Return(nil, errors.New("get error"))

	tokens, err := s.service.RefreshTokens(s.ctx, oldRefresh)

	s.Error(err)
	s.Nil(tokens)
	s.refreshRepo.AssertExpectations(s.T())
}

// --- UpdatePassword ---

func (s *AuthServiceSuite) TestUpdatePassword_Success() {
	userID := s.objMother.DefaultUserID()
	passwords := s.objMother.DefaultPasswords()
	hashed := s.objMother.DefaultHashedPassword()
	newHashed := "new_hashed"

	s.authRepo.On("GetPasswordByID", s.ctx, userID).Return(hashed, nil)
	s.hasher.On("CompareHashAndPassword", hashed, passwords.Old).Return(nil)
	s.hasher.On("GenerateFromPassword", passwords.New).Return(newHashed, nil)
	s.authRepo.On("UpdatePassword", s.ctx, userID, newHashed).Return(nil)

	err := s.service.UpdatePassword(s.ctx, userID, passwords)
	s.NoError(err)
	s.authRepo.AssertExpectations(s.T())
	s.hasher.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestUpdatePassword_GetPasswordError() {
	userID := s.objMother.DefaultUserID()
	s.authRepo.On("GetPasswordByID", s.ctx, userID).Return("", errors.New("db error"))

	err := s.service.UpdatePassword(s.ctx, userID, &entity.UserPasswords{Old: "old", New: "new"})
	s.Error(err)
	s.authRepo.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestUpdatePassword_InvalidOldPassword() {
	userID := s.objMother.DefaultUserID()
	hashed := s.objMother.DefaultHashedPassword()
	s.authRepo.On("GetPasswordByID", s.ctx, userID).Return(hashed, nil)
	s.hasher.On("CompareHashAndPassword", hashed, "wrong").Return(errors.New("mismatch"))

	err := s.service.UpdatePassword(s.ctx, userID, &entity.UserPasswords{Old: "wrong", New: "new"})

	s.Error(err)
	s.ErrorIs(err, ErrInvalidCredentials)
	s.authRepo.AssertExpectations(s.T())
	s.hasher.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestGetClaims_Success() {
	claims := s.objMother.DefaultClaims()
	token := "access"
	s.tokenService.On("ParseAccessToken", token).Return(claims, nil)

	parsed, err := s.service.GetClaims(s.ctx, token)
	s.NoError(err)
	s.Equal(claims, parsed)

	s.tokenService.AssertExpectations(s.T())
}

func (s *AuthServiceSuite) TestGetClaims_ParseError() {
	token := "badtoken"
	s.tokenService.On("ParseAccessToken", token).Return(nil, errors.New("parse error"))

	claims, err := s.service.GetClaims(s.ctx, token)
	s.Error(err)
	s.Nil(claims)

	s.tokenService.AssertExpectations(s.T())
}
