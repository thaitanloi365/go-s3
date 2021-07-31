package s3

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type UploadFileParams struct {
	Data               io.Reader
	Bucket             string
	Key                string
	Metadata           map[string]*string
	ContentDisposition string
	ContentType        string
	ACL                string
	CacheControl       *string
	CacheExpires       *time.Time
}

func (client *Client) UploadFile(params UploadFileParams) (string, error) {
	s, err := client.NewSession()
	if err != nil {
		client.logger.Printf("Create session err: %v\n", err)
		return "", err
	}

	var svc = s3manager.NewUploader(s)

	var uploadParams = &s3manager.UploadInput{
		Bucket:       aws.String(params.Bucket),
		Key:          aws.String(params.Key),
		ACL:          aws.String("public-read"),
		Body:         params.Data,
		CacheControl: params.CacheControl,
		Expires:      params.CacheExpires,
	}

	if params.ACL != "" {
		uploadParams.ACL = aws.String(params.ACL)
	}

	if params.ContentDisposition != "" {
		uploadParams.ContentDisposition = aws.String(params.ContentDisposition)
	}

	if params.ContentType != "" {
		uploadParams.ContentType = aws.String(params.ContentType)
	}

	if params.Metadata != nil {
		uploadParams.Metadata = params.Metadata
	}

	result, err := svc.Upload(uploadParams)
	if err != nil {
		client.logger.Printf("Upload s3 err: %v\n", err)
		return "", err
	}

	return result.Location, nil

}

func (client *Client) UploadFiles(params []UploadFileParams) (result []string) {
	var wg sync.WaitGroup
	var max = len(params)
	var urlChannel = make(chan string, max)
	wg.Add(max)

	for _, v := range params {
		go func(wg *sync.WaitGroup, param UploadFileParams, channel chan string) {
			defer func() {
				wg.Done()
				max--
				if max == 0 {
					close(channel)
				}

			}()

			resp, _ := client.UploadFile(param)
			channel <- resp
		}(&wg, v, urlChannel)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			value, more := <-urlChannel
			if more {
				if value != "" {
					result = append(result, value)
				}
			} else {
				return
			}
		}
	}()

	wg.Wait()

	return
}

type UploadLogParams struct {
	IgnoreFiles               []string
	FolderToUpload            string
	UploadToBucket            string
	ShouldKeepFileAfterUpload func(filePath, fileName string) bool
}

func (client *Client) UploadLog(params UploadLogParams) ([]string, error) {
	var response = []string{}
	var files []string
	var dir = params.FolderToUpload

	var ignoreFiles = params.IgnoreFiles
	var err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Size() > 0 {
			if len(ignoreFiles) > 0 {
				for _, ignoreFile := range ignoreFiles {
					if !strings.Contains(path, ignoreFile) {
						files = append(files, path)
					}
				}
			} else {
				files = append(files, path)
			}
		}

		return nil
	})
	if err != nil {
		return response, err
	}

	if len(files) == 0 {
		client.logger.Printf("No have any files to upload\n")
		return response, nil
	}

	var wg sync.WaitGroup
	var max = len(files)
	var urlChannel = make(chan string, max)
	wg.Add(max)

	for _, file := range files {
		go func(wg *sync.WaitGroup, file string, channel chan string) {
			defer func() {
				wg.Done()
				max--
				if max == 0 {
					close(channel)
				}

			}()

			originFile, err := os.Open(file)
			if err != nil {
				client.logger.Printf("Open file: error %+v\n", err)
				channel <- ""
				return
			}

			reader, writer := io.Pipe()
			go func() {
				gw := gzip.NewWriter(writer)
				io.Copy(gw, originFile)
				originFile.Close()
				gw.Close()
				writer.Close()
			}()
			var ext = path.Ext(file)
			var fileName = file[0 : len(file)-len(ext)]
			var gzipFileName = fmt.Sprintf("%s.gz", fileName)
			var fileKey = filepath.Base(gzipFileName)
			var folder = filepath.Dir(file)
			var key = fmt.Sprintf("%s/%s", folder, fileKey)

			result, _ := client.UploadFile(UploadFileParams{
				Data:   reader,
				Bucket: params.UploadToBucket,
				Key:    key,
			})
			channel <- result

			if result != "" {
				if params.ShouldKeepFileAfterUpload != nil && params.ShouldKeepFileAfterUpload(key, fileName) {
					os.Remove(file)

				}
			}
		}(&wg, file, urlChannel)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			value, more := <-urlChannel
			if more {
				if value != "" {
					response = append(response, value)
				}
			} else {
				return
			}
		}
	}()

	wg.Wait()

	return response, nil
}
