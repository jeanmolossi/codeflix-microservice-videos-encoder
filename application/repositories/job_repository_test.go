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

func TestJobRepositoryDb_Insert(t *testing.T) {
	Db := database.NewDbTest()
	defer Db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "fake-path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: Db}
	repo.Insert(video)

	job, err := domain.NewJob("output-path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: Db}
	repoJob.Insert(job)

	j, err := repoJob.Find(job.ID)

	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.ID, job.ID)
	require.Equal(t, j.VideoID, video.ID)
}

func TestJobRepositoryDb_Update(t *testing.T) {
	Db := database.NewDbTest()
	defer Db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "fake-path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: Db}
	repo.Insert(video)

	job, err := domain.NewJob("output-path", "Pending", video)
	require.Nil(t, err)

	repoJob := repositories.JobRepositoryDb{Db: Db}

	job.Status = "Complete"

	repoJob.Update(job)

	j, err := repoJob.Find(job.ID)
	require.NotEmpty(t, j.ID)
	require.Nil(t, err)
	require.Equal(t, j.Status, job.Status)
}
