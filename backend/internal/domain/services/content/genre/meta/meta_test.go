package genre_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	genre "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/genre/meta"
	usecase "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/content/genre"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
)

type GenreObjectMother struct{}

func (GenreObjectMother) AdminClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.Admin}
}

func (GenreObjectMother) UserClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.User}
}

func (GenreObjectMother) ValidGenre() *entity.Genre {
	return &entity.Genre{ID: uuid.New(), Title: "Rock"}
}

type GenreServiceSuite struct {
	suite.Suite
	ctx    context.Context
	repo   *mocks.GenreRepository
	svc    *genre.GenreService
	mother GenreObjectMother
}

func TestGenreServiceSuite(t *testing.T) {
	suite.Run(t, new(GenreServiceSuite))
}

func (s *GenreServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewGenreRepository(s.T())
	s.svc = genre.NewGenreService(s.repo)
	s.mother = GenreObjectMother{}
}

func (s *GenreServiceSuite) TearDownTest() {
	s.repo.AssertExpectations(s.T())
}

func (s *GenreServiceSuite) TestCreateGenre() {
	admin := s.mother.AdminClaims()
	user := s.mother.UserClaims()
	valid := s.mother.ValidGenre()

	s.Run("success", func() {
		s.repo.On("Create", s.ctx, valid).Return(nil)
		err := s.svc.CreateGenre(s.ctx, admin, valid)
		s.NoError(err)
	})

	s.Run("forbidden", func() {
		err := s.svc.CreateGenre(s.ctx, user, valid)
		s.ErrorIs(err, commonerr.ErrForbidden)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.repo.On("Create", s.ctx, valid).Return(errors.New("db error"))
		err := s.svc.CreateGenre(s.ctx, admin, valid)
		s.ErrorIs(err, usecase.ErrCreateGenre)
	})
}

func (s *GenreServiceSuite) TestGetGenreByID() {
	valid := s.mother.ValidGenre()

	s.Run("success", func() {
		s.SetupTest()
		s.repo.On("GetByID", s.ctx, valid.ID).Return(valid, nil)
		res, err := s.svc.GetGenreByID(s.ctx, valid.ID)
		s.NoError(err)
		s.Equal(valid, res)
	})

	s.Run("invalid ID", func() {
		s.SetupTest()
		_, err := s.svc.GetGenreByID(s.ctx, uuid.Nil)
		s.ErrorIs(err, genre.ErrInvalidGenreID)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.repo.On("GetByID", s.ctx, valid.ID).Return(nil, errors.New("repo error"))
		_, err := s.svc.GetGenreByID(s.ctx, valid.ID)
		s.ErrorIs(err, usecase.ErrGetGenre)
	})
}

func (s *GenreServiceSuite) TestUpdateGenre() {
	admin := s.mother.AdminClaims()
	user := s.mother.UserClaims()
	valid := s.mother.ValidGenre()

	s.Run("success", func() {
		s.SetupTest()
		s.repo.On("Update", s.ctx, valid).Return(nil)
		err := s.svc.UpdateGenre(s.ctx, admin, valid)
		s.NoError(err)
	})

	s.Run("forbidden", func() {
		s.SetupTest()
		err := s.svc.UpdateGenre(s.ctx, user, valid)
		s.ErrorIs(err, commonerr.ErrForbidden)
	})

	s.Run("invalid ID", func() {
		s.SetupTest()
		err := s.svc.UpdateGenre(s.ctx, admin, &entity.Genre{ID: uuid.Nil})
		s.ErrorIs(err, genre.ErrInvalidGenreID)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.repo.On("Update", s.ctx, valid).Return(errors.New("db error"))
		err := s.svc.UpdateGenre(s.ctx, admin, valid)
		s.ErrorIs(err, usecase.ErrUpdateGenre)
	})
}

func (s *GenreServiceSuite) TestDeleteGenre() {
	admin := s.mother.AdminClaims()
	user := s.mother.UserClaims()
	id := uuid.New()

	s.Run("success", func() {
		s.SetupTest()
		s.repo.On("Delete", s.ctx, id).Return(nil)
		err := s.svc.DeleteGenre(s.ctx, admin, id)
		s.NoError(err)
	})

	s.Run("forbidden", func() {
		s.SetupTest()
		err := s.svc.DeleteGenre(s.ctx, user, id)
		s.ErrorIs(err, commonerr.ErrForbidden)
	})

	s.Run("invalid ID", func() {
		s.SetupTest()
		err := s.svc.DeleteGenre(s.ctx, admin, uuid.Nil)
		s.ErrorIs(err, genre.ErrInvalidGenreID)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		genreID := uuid.New()
		repo := mocks.NewGenreRepository(s.T())
		repo.On("Delete", s.ctx, genreID).Return(errors.New("db error"))

		svc := genre.NewGenreService(repo)
		err := svc.DeleteGenre(s.ctx, admin, genreID)

		assert.Error(s.T(), err)
		assert.True(s.T(), errors.Is(err, usecase.ErrDeleteGenre))
	})
}

func (s *GenreServiceSuite) TestGetGenreByAlbum() {
	albumID := uuid.New()
	genres := []*entity.Genre{s.mother.ValidGenre()}

	s.Run("success", func() {
		s.SetupTest()
		s.repo.On("GetByAlbum", s.ctx, albumID).Return(genres, nil)
		res, err := s.svc.GetGenreByAlbum(s.ctx, albumID)
		s.NoError(err)
		s.Equal(genres, res)
	})

	s.Run("repo error", func() {
		s.SetupTest()
		s.repo.On("GetByAlbum", s.ctx, albumID).Return(nil, errors.New("db error"))
		res, err := s.svc.GetGenreByAlbum(s.ctx, albumID)
		s.Nil(res)
		s.Error(err)
		s.True(errors.Is(err, usecase.ErrGetGenreByAlbum))
	})
}
