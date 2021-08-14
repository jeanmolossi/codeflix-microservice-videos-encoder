package repositories_test

import (
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/application/repositories"
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/domain"
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestVideoRepositoryDb_Insert(t *testing.T) {
	Db := database.NewDbTest()
	defer Db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "fake-path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: Db}
	repo.Insert(video)

	v, err := repo.Find(video.ID)

	require.NotEmpty(t, v.ID)
	require.Nil(t, err)
	require.Equal(t, v.ID, video.ID)
}
