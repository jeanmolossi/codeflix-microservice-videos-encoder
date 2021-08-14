package domain_test

import (
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVideoValidateIfVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()

	err := video.Validate()

	require.Error(t, err)
}

func TestVideoIdIsNotUuid(t *testing.T) {
	video := domain.NewVideo()

	video.ID = "not-uuid"
	video.ResourceID = "not-uuid"
	video.FilePath = "fake-path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Error(t, err)
}

func TestVideoValidation(t *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.ResourceID = "not-null-value"
	video.FilePath = "fake-path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Nil(t, err)
}
