package s3

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadFile(t *testing.T) {
	file, err := os.Open("go.png")
	assert.NoError(t, err)

	url, err := New(&Config{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
		Region:    os.Getenv("AWS_REGION"),
	}).UploadFile(UploadFileParams{
		Data:               file,
		Bucket:             os.Getenv("AWS_BUCKET"),
		Key:                "test.png",
		ContentType:        "image/jpeg",
		ContentDisposition: "inline",
	})
	assert.NoError(t, err)
	assert.NotEqual(t, url, "")

	fmt.Println("url", url)
}

func TestUploadFiles(t *testing.T) {
	file, err := os.Open("go.png")
	assert.NoError(t, err)

	var urls = New(&Config{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
		Region:    os.Getenv("AWS_REGION"),
	}).UploadFiles([]UploadFileParams{
		{
			Data:               file,
			Bucket:             os.Getenv("AWS_BUCKET"),
			Key:                "test1.png",
			ContentType:        "image/jpeg",
			ContentDisposition: "inline",
		},
		{
			Data:               file,
			Bucket:             os.Getenv("AWS_BUCKET"),
			Key:                "test2.png",
			ContentType:        "image/jpeg",
			ContentDisposition: "inline",
		},
	})
	assert.NoError(t, err)
	assert.NotEqual(t, urls, []string{})

	fmt.Println("urls", urls)
}

func TestUploadLog(t *testing.T) {
	urls, err := New(&Config{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
		Region:    os.Getenv("AWS_REGION"),
	}).UploadLog(UploadLogParams{
		FolderToUpload: "logs",
		UploadToBucket: os.Getenv("AWS_BUCKET"),
	})
	assert.NoError(t, err)
	assert.NotEqual(t, urls, []string{})

	fmt.Println("urls", urls)
}
