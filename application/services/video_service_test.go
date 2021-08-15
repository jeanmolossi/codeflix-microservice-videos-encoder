package services_test

import (
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/application/repositories"
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/application/services"
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/domain"
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/framework/database"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func getRandomVideo() string {
	var videos []string
	videos = append(
		videos,
		"BigBuckBunny.mp4",
		"ForBiggerFun.mp4",
		"TearsOfSteel.mp4",
	)

	return videos[rand.Intn(len(videos))]
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	Db := database.NewDbTest()
	defer Db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = getRandomVideo()
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: Db}

	return video, repo
}

func TestVideoService_Download(t *testing.T) {
	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download(os.Getenv("AWS_S3_BUCKET"))
	require.Nil(t, err)

	err = videoService.Fragment()
	require.Nil(t, err)

	err = videoService.Encode()
	require.Nil(t, err)

	err = videoService.Finish()
	require.Nil(t, err)
}
