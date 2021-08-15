package services_test

import (
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/application/services"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestVideoUpload_ProcessUpload(t *testing.T) {
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

	videoUpload := services.NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("AWS_S3_BUCKET")
	videoUpload.VideoPath = os.Getenv("LOCAL_STORAGE_PATH") + "/" + video.ID

	doneUpload := make(chan string)
	go videoUpload.ProcessUpload(50, doneUpload)

	result := <-doneUpload
	require.Equal(t, result, "upload completed")

	err = videoService.Finish()
	require.Nil(t, err)
}
