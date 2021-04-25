package s3

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type Client struct {
	config *Config
	logger Logger
	Logger Logger
}

type Config struct {
	AccessKey string
	SecretKey string
	Region    string
	Logger    Logger
}

func New(config *Config) *Client {
	var client = &Client{
		config: config,
		logger: log.New(os.Stdout, "\r\n", 0),
	}
	if config.Logger != nil {
		client.logger = config.Logger
	}

	return client
}

func (client *Client) NewSession() (*session.Session, error) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(client.config.Region),
		Credentials: credentials.NewStaticCredentials(client.config.AccessKey, client.config.SecretKey, ""),
	})

	return session, err
}
