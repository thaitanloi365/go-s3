package s3

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
)

func TestDownloadFile(t *testing.T) {

	var params = DownloadFileParams{
		Key:    "orders/bvm06gp0e7223fc0rah0/invoices/4a9cba59e000005",
		Bucket: os.Getenv("AWS_BUCKET"),
		Output: aws.NewWriteAtBuffer(nil),
	}

	err := New(&Config{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
		Region:    os.Getenv("AWS_REGION"),
	}).DownloadFile(&params)
	assert.NoError(t, err)

	fmt.Println(len(params.Output.Bytes()))

	ioutil.WriteFile("4a9cba59e000005.pdf", params.Output.Bytes(), 0)
}

func TestDownloadFiles(t *testing.T) {
	var params = []*DownloadFileParams{
		{
			Key:      "orders/bvm06gp0e7223fc0rah0/invoices/4a9cba59e000005",
			Bucket:   os.Getenv("AWS_BUCKET"),
			Output:   aws.NewWriteAtBuffer(nil),
			FileName: "4a9cba59e000005.pdf",
		},
		{
			Key:      "orders/bubsu7fi3h5pjjf992p0/invoices/4dc4b0b1300000a",
			Bucket:   os.Getenv("AWS_BUCKET"),
			Output:   aws.NewWriteAtBuffer(nil),
			FileName: "4dc4b0b1300000a.pdf",
		},
	}

	err := New(&Config{
		AccessKey: os.Getenv("AWS_ACCESS_KEY"),
		SecretKey: os.Getenv("AWS_SECRET_KEY"),
		Region:    os.Getenv("AWS_REGION"),
	}).DownloadFiles(params)
	assert.NoError(t, err)

	for _, param := range params {
		ioutil.WriteFile(param.FileName, param.Output.Bytes(), 0)
	}

}
