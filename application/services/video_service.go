package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/application/repositories"
	"github.com/jeanmolossi/codeflix-microservice-videos-encoder/domain"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func getCredentials() (aws.Credentials, error) {
	if os.Getenv("AWS_ACCESS_KEY") == "" {
		return aws.Credentials{}, fmt.Errorf("aws_access_key is empty: %s", os.Getenv("AWS_ACCESS_KEY"))
	}

	if os.Getenv("AWS_SECRET_KEY") == "" {
		return aws.Credentials{}, fmt.Errorf("aws_secret_key is empty: %s", os.Getenv("AWS_SECRET_KEY"))
	}

	awsCredentials := aws.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY"),
		SecretAccessKey: os.Getenv("AWS_SECRET_KEY"),
		Source:          "Env file, Hard coded",
	}

	return awsCredentials, nil
}

func (v *VideoService) Download(bucketName string) error {
	sourceBucket := bucketName
	objectName := aws.String("ForBiggerFun.mp4")

	if sourceBucket == "" || *objectName == "" {
		return fmt.Errorf("you must supply the bucket to copy %v", sourceBucket)
	}

	awsCredentials, err := getCredentials()
	if err != nil {
		log.Fatalf("Error on load awsCrededntials: %v", err)
		return err
	}

	ctx := context.Background()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	awsConfig.Region = "us-east-1"
	awsConfig.Credentials = credentials.StaticCredentialsProvider{
		Value: awsCredentials,
	}

	if err != nil {
		fmt.Printf("configuration error, %v", err.Error())
		return err
	}

	client := s3.NewFromConfig(awsConfig)

	objectInput := &s3.GetObjectInput{
		Bucket: aws.String(url.PathEscape(sourceBucket)),
		Key:    objectName,
	}

	objOutput, err := client.GetObject(ctx, objectInput)

	file, err := ioutil.ReadAll(objOutput.Body)

	f, err := os.Create(v.getFile())
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(file)
	if err != nil {
		return err
	}

	log.Printf("Video has been stored")

	return nil
}

func (v *VideoService) Fragment() error {
	err := os.Mkdir(v.getFolder(), os.ModePerm)
	if err != nil {
		return err
	}

	source := v.getFile()
	target := v.getFragFile()

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Encode() error {
	var cmdArgs []string
	cmdArgs = append(cmdArgs, v.getFragFile())
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, v.getFolder())
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")
	cmd := exec.Command("mp4dash", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return err
	}

	printOutput(output)

	return nil
}

func (v *VideoService) Finish() error {
	err := os.Remove(v.getFile())
	if err != nil {
		log.Println("error removing mp4 ", v.Video.ID, ".mp4")
		return err
	}

	err = os.Remove(v.getFragFile())
	if err != nil {
		log.Println("error removing frag ", v.Video.ID, ".frag")
		return err
	}

	err = os.RemoveAll(v.getFolder())
	if err != nil {
		log.Println("error removing folder", v.Video.ID, "folder")
		return err
	}

	log.Println("files has been removed: ", v.Video.ID)
	return nil
}

func (v *VideoService) getFolder() string {
	return os.Getenv("LOCAL_STORAGE_PATH") + "/" + v.Video.ID
}

func (v *VideoService) getFile() string {
	return v.getFolder() + ".mp4"
}

func (v *VideoService) getFragFile() string {
	return v.getFolder() + ".frag"
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n", string(out))
	}
}
