package aws

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// ClientConfig is the client's config
type ClientConfig struct {
	Region  string
	RoleARN string
}

var _ FileSystem = &Client{}

// Client is the client
type Client struct {
	Config *ClientConfig
}

// Get gets the value param from ssm
func (c *Client) Get(pattern string) (string, error) {
	client := c.ssm()

	params := &ssm.GetParameterInput{
		Name:           aws.String(pattern),
		WithDecryption: aws.Bool(true),
	}

	response, err := client.GetParameter(params)
	if err != nil {
		return "", err
	}

	if param := response.Parameter; param != nil {
		if value := param.Value; value != nil {
			return *value, nil
		}
	}

	return "", nil
}

// Glob returns a list of all paths
func (c *Client) Glob(bucket, pattern string) ([]string, error) {
	client := c.s3()

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
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
func (c *Client) ReadFile(bucket, path string) ([]byte, error) {
	client := c.s3()

	params := &s3.GetObjectInput{
		Key:    aws.String(path),
		Bucket: aws.String(bucket),
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

func (c *Client) s3() *s3.S3 {
	config := &aws.Config{
		Region: aws.String(c.Config.Region),
	}

	cookie := session.Must(session.NewSession(config))

	if c.Config.RoleARN != "" {
		creds := stscreds.NewCredentials(cookie, c.Config.RoleARN)
		return s3.New(cookie, &aws.Config{Credentials: creds})
	}

	return s3.New(cookie)
}

func (c *Client) ssm() *ssm.SSM {
	config := &aws.Config{
		Region: aws.String(c.Config.Region),
	}

	cookie := session.Must(session.NewSession(config))

	if c.Config.RoleARN != "" {
		creds := stscreds.NewCredentials(cookie, c.Config.RoleARN)
		return ssm.New(cookie, &aws.Config{Credentials: creds})
	}

	return ssm.New(session.Must(session.NewSession(config)))
}
