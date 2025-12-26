package audio_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hahaclassic/orpheon/backend/internal/domain/entity"
	audio "github.com/hahaclassic/orpheon/backend/internal/domain/services/content/track/audio"
	commonerr "github.com/hahaclassic/orpheon/backend/internal/domain/usecases/errors"
	"github.com/hahaclassic/orpheon/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type AudioFileServiceSuite struct {
	suite.Suite
	service   *audio.AudioFileService
	repo      *mocks.AudioFileRepository
	converter *mocks.AudioConverter
	ctx       context.Context
	trackID   uuid.UUID
}

func TestAudioFileServiceSuite(t *testing.T) {
	suite.Run(t, new(AudioFileServiceSuite))
}

func (s *AudioFileServiceSuite) SetupTest() {
	s.repo = mocks.NewAudioFileRepository(s.T())
	s.converter = mocks.NewAudioConverter(s.T())
	s.service = audio.New(s.repo, s.converter)
	s.ctx = context.Background()
	s.trackID = uuid.New()
}

// Object Mother
func ValidAudioChunk(trackID uuid.UUID, size int64) *entity.AudioChunk {
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i)
	}
	return &entity.AudioChunk{
		TrackID: trackID,
		Start:   0,
		End:     size,
		Data:    data,
	}
}

func InvalidAudioChunk() *entity.AudioChunk {
	return &entity.AudioChunk{
		TrackID: uuid.New(),
		Start:   10,
		End:     5,
		Data:    []byte{},
	}
}

func AdminClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.Admin}
}

func UserClaims() *entity.Claims {
	return &entity.Claims{AccessLvl: entity.User}
}

// GetAudioChunk
func (s *AudioFileServiceSuite) TestGetAudioChunkValid() {
	chunk := ValidAudioChunk(s.trackID, 10)
	s.repo.On("GetAudioChunk", mock.Anything, chunk).Return(chunk, nil)

	res, err := s.service.GetAudioChunk(s.ctx, chunk)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), chunk, res)
}

func (s *AudioFileServiceSuite) TestGetAudioChunkInvalidParams() {
	chunk := InvalidAudioChunk()
	res, err := s.service.GetAudioChunk(s.ctx, chunk)
	assert.ErrorIs(s.T(), err, audio.ErrInvalidChunkParams)
	assert.Nil(s.T(), res)
}

func (s *AudioFileServiceSuite) TestGetAudioChunkRepoError() {
	chunk := ValidAudioChunk(s.trackID, 10)
	s.repo.On("GetAudioChunk", mock.Anything, chunk).Return(nil, errors.New("repo error"))

	res, err := s.service.GetAudioChunk(s.ctx, chunk)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), res)
}

// UploadAudioFile
func (s *AudioFileServiceSuite) TestUploadAudioFileValid() {
	chunk := ValidAudioChunk(s.trackID, 10)
	s.converter.On("ChangeBitrate", mock.Anything, chunk).Return(chunk, nil)
	s.repo.On("UploadAudioFile", mock.Anything, chunk).Return(nil)

	err := s.service.UploadAudioFile(s.ctx, AdminClaims(), chunk)
	assert.NoError(s.T(), err)
}

func (s *AudioFileServiceSuite) TestUploadAudioFileNotAdmin() {
	chunk := ValidAudioChunk(s.trackID, 10)
	err := s.service.UploadAudioFile(s.ctx, UserClaims(), chunk)
	assert.ErrorIs(s.T(), err, commonerr.ErrForbidden)
}

func (s *AudioFileServiceSuite) TestUploadAudioFileInvalidChunk() {
	chunk := InvalidAudioChunk()
	err := s.service.UploadAudioFile(s.ctx, AdminClaims(), chunk)
	assert.ErrorIs(s.T(), err, audio.ErrInvalidChunkParams)
}

// DeleteAudioFile
func (s *AudioFileServiceSuite) TestDeleteAudioFileValid() {
	s.repo.On("DeleteFile", mock.Anything, s.trackID).Return(nil)
	err := s.service.DeleteAudioFile(s.ctx, AdminClaims(), s.trackID)
	assert.NoError(s.T(), err)
}

func (s *AudioFileServiceSuite) TestDeleteAudioFileNotAdmin() {
	err := s.service.DeleteAudioFile(s.ctx, UserClaims(), s.trackID)
	assert.ErrorIs(s.T(), err, commonerr.ErrForbidden)
}

func (s *AudioFileServiceSuite) TestDeleteAudioFileRepoError() {
	s.repo.On("DeleteFile", mock.Anything, s.trackID).Return(errors.New("delete error"))
	err := s.service.DeleteAudioFile(s.ctx, AdminClaims(), s.trackID)
	assert.Error(s.T(), err)
}
