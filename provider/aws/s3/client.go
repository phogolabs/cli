package s3

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// ClientConfig is the client's config
type ClientConfig struct {
	Region  string
	Bucket  string
	RoleARN string
}

var _ FileSystem = &Client{}

// Client is the client
type Client struct {
	Config *ClientConfig
}

// Glob returns a list of all paths
func (c *Client) Glob(pattern string) ([]string, error) {
	client := c.client()

	params := &s3.ListObjectsInput{
		Bucket: aws.String(c.Config.Bucket),
	}

	response, err := client.ListObjects(params)
	if err != nil {
		return nil, err
	}

	var paths []string

	for _, item := range response.Contents {
		matched, err := filepath.Match(pattern, *item.Key)

		switch {
		case err != nil:
			return nil, err
		case matched:
			paths = append(paths, *item.Key)
		}
	}

	return paths, nil
}

// ReadFile reads a file from the bucket
func (c *Client) ReadFile(path string) ([]byte, error) {
	client := c.client()

	params := &s3.GetObjectInput{
		Key:    aws.String(path),
		Bucket: aws.String(c.Config.Bucket),
	}

	response, err := client.GetObject(params)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	buffer := &bytes.Buffer{}
	if _, err := io.Copy(buffer, response.Body); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (c *Client) client() *s3.S3 {
	config := &aws.Config{
		Region: aws.String(c.Config.Region),
	}

	cookie := session.New(config)

	if c.Config.RoleARN != "" {
		creds := stscreds.NewCredentials(cookie, c.Config.RoleARN)
		return s3.New(cookie, &aws.Config{Credentials: creds})
	}

	return s3.New(cookie)
}
