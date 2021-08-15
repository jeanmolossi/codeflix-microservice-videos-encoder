package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectPath string, client *s3.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("LOCAL_STORAGE_PATH")+"/")

	f, err := os.Open(objectPath)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String(path[1]),
		Body:   f,
	})

	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	err := vu.loadPaths()
	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()
	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {
		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

func (vu *VideoUpload) uploadWorker(in chan int, returnChan chan string, uploadClient *s3.Client, ctx context.Context) {

	for x := range in {
		err := vu.UploadObject(vu.Paths[x], uploadClient, ctx)

		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[x])
			log.Printf("error during upload: %v, error: %v", vu.Paths[x], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}

	returnChan <- "upload completed"
}

func getClientUpload() (*s3.Client, context.Context, error) {
	ctx := context.Background()

	awsCredentials, err := getCredentials()
	if err != nil {
		log.Fatalf("Error on load awsCrededntials: %v", err)
		return nil, ctx, err
	}

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, ctx, err
	}

	awsConfig.Region = "us-east-1"
	awsConfig.Credentials = credentials.StaticCredentialsProvider{
		Value: awsCredentials,
	}

	if err != nil {
		fmt.Printf("configuration error, %v", err.Error())
		return nil, ctx, err
	}

	client := s3.NewFromConfig(awsConfig)

	return client, ctx, nil
}
