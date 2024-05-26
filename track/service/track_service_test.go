package service

import (
	"testing"

	"github.com/rauzh/cd-core/models"

	mocks "github.com/rauzh/cd-core/repo/mocks"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
)

func TestTrackService_Get(t *testing.T) {

	mockTrackRepo := mocks.NewTrackRepo(t)
	mockTrackRepo.EXPECT().Get(mock.AnythingOfType("context.backgroundCtx"), uint64(1234)).Return(&models.Track{
		TrackID:  1234,
		Title:    "aa",
		Duration: 120,
		Genre:    "rock",
		Artists:  []uint64{82, 4},
	}, nil).Once()

	ts := NewTrackService(mockTrackRepo)

	track, err := ts.Get(1234)
	assert.Nil(t, err)
	assert.Equal(t, "rock", track.Genre)
}