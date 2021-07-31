package s3

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	ErrNotFound = errors.New("not found")
)

func (client *Client) CheckFile(bucket, key string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		client.logger.Printf("Create session err: %v\n", err)
		return "", err
	}

	_, err = s3.New(session).HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return "", ErrNotFound
			default:
				return "", err
			}
		}
		return "", err
	}

	var url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)

	return url, nil
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}

	return false
}
