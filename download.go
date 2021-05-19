package s3

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type DownloadFileParams struct {
	Bucket   string
	Key      string
	Output   *aws.WriteAtBuffer
	FileName string
}

func (client *Client) DownloadFile(params *DownloadFileParams) error {
	session, err := client.NewSession()
	if err != nil {
		client.logger.Printf("Create session err: %v\n", err)
		return err
	}
	var downloader = s3manager.NewDownloaderWithClient(s3.New(session))
	_, err = downloader.Download(params.Output, &s3.GetObjectInput{
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(params.Key),
	})
	if err != nil {
		client.logger.Printf("Download file from bucket = %s key = %s error: %v\n", params.Bucket, params.Key, err)
		return err
	}

	return nil
}

func (client *Client) DownloadFiles(params []*DownloadFileParams) error {
	var wg = sync.WaitGroup{}

	for _, param := range params {
		wg.Add(1)
		go func(param *DownloadFileParams) {
			defer wg.Done()
			client.DownloadFile(param)

		}(param)
	}

	wg.Wait()

	return nil
}
