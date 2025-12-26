package user_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	usersvc "github.com/hahaclassic/orpheon/backend/internal/domain/services/user"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, &UserServiceSuite{})
}

type UserObjectMother struct{}

func (UserObjectMother) DefaultUser() *entity.UserInfo {
	return &entity.UserInfo{Name: "testuser"}
}

func (UserObjectMother) DefaultClaims(userID uuid.UUID, accessLvl entity.AccessLevel) *entity.Claims {
	return &entity.Claims{UserID: userID, AccessLvl: accessLvl}
}

type UserServiceSuite struct {
	suite.Suite

	ctx     context.Context
	service *usersvc.UserService
	repo    *mocks.UserRepository

	objMother *UserObjectMother
}

func (s *UserServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewUserRepository(s.T())
	s.service = usersvc.New(s.repo)
	s.objMother = &UserObjectMother{}
}

// --- CreateUser ---

func (s *UserServiceSuite) TestCreateUser_Success() {
	user := s.objMother.DefaultUser()
	s.repo.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	id, err := s.service.CreateUser(s.ctx, user)

	s.NoError(err)
	s.NotEqual(uuid.Nil, id)
	s.Equal(entity.User, user.AccessLvl)
	s.WithinDuration(time.Now(), user.RegistrationDate, time.Second)
	s.repo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestCreateUser_RepoError() {
	user := s.objMother.DefaultUser()
	s.repo.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("db error"))

	id, err := s.service.CreateUser(s.ctx, user)

	s.Error(err)
	s.Equal(uuid.Nil, id)
	s.repo.AssertExpectations(s.T())
}

// --- GetUser ---

func (s *UserServiceSuite) TestGetUser_Success() {
	userID := uuid.New()
	expected := &entity.UserInfo{ID: userID, Name: "testuser"}
	s.repo.On("GetUser", mock.Anything, userID).Return(expected, nil)

	result, err := s.service.GetUser(s.ctx, userID)

	s.NoError(err)
	s.Equal(expected, result)
	s.repo.AssertExpectations(s.T())
}

// --- UpdateUser ---

func (s *UserServiceSuite) TestUpdateUser_Success() {
	userID := uuid.New()
	user := &entity.UserInfo{ID: userID}
	claims := s.objMother.DefaultClaims(userID, entity.User)
	s.repo.On("UpdateUser", mock.Anything, user).Return(nil)

	err := s.service.UpdateUser(s.ctx, claims, user)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestUpdateUser_Forbidden() {
	s.SetupTest()
	user := &entity.UserInfo{ID: uuid.New()}
	claims := s.objMother.DefaultClaims(uuid.New(), entity.User)

	err := s.service.UpdateUser(s.ctx, claims, user)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
}

func (s *UserServiceSuite) TestUpdateUser_RepoError() {
	userID := uuid.New()
	user := &entity.UserInfo{ID: userID}
	claims := s.objMother.DefaultClaims(userID, entity.User)
	s.repo.On("UpdateUser", mock.Anything, user).Return(errors.New("db error"))

	err := s.service.UpdateUser(s.ctx, claims, user)

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

// --- DeleteUser ---

func (s *UserServiceSuite) TestDeleteUser_Self() {
	userID := uuid.New()
	claims := s.objMother.DefaultClaims(userID, entity.User)
	s.repo.On("DeleteUser", mock.Anything, userID).Return(nil)

	err := s.service.DeleteUser(s.ctx, claims, userID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestDeleteUser_Admin() {
	adminID := uuid.New()
	userID := uuid.New()
	claims := s.objMother.DefaultClaims(adminID, entity.Admin)
	s.repo.On("DeleteUser", mock.Anything, userID).Return(nil)

	err := s.service.DeleteUser(s.ctx, claims, userID)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *UserServiceSuite) TestDeleteUser_Forbidden() {
	s.SetupTest()
	userID := uuid.New()
	claims := s.objMother.DefaultClaims(uuid.New(), entity.User)

	err := s.service.DeleteUser(s.ctx, claims, userID)

	s.Error(err)
	s.ErrorIs(err, commonerr.ErrForbidden)
}

func (s *UserServiceSuite) TestDeleteUser_RepoError() {
	userID := uuid.New()
	claims := s.objMother.DefaultClaims(userID, entity.User)
	s.repo.On("DeleteUser", mock.Anything, userID).Return(errors.New("db error"))

	err := s.service.DeleteUser(s.ctx, claims, userID)

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}
