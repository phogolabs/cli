package ssm

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// ClientConfig is the client's config
type ClientConfig struct {
	Region  string
	RoleARN string
}

var _ Getter = &Client{}

// Client is the client
type Client struct {
	Config *ClientConfig
}

// Get gets the value param from ssm
func (c *Client) Get(pattern string) (string, error) {
	client := c.client()

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

func (c *Client) client() *ssm.SSM {
	config := &aws.Config{
		Region: aws.String(c.Config.Region),
	}

	cookie := session.New(config)

	if c.Config.RoleARN != "" {
		creds := stscreds.NewCredentials(cookie, c.Config.RoleARN)
		return ssm.New(cookie, &aws.Config{Credentials: creds})
	}

	return ssm.New(session.New(config))
}
