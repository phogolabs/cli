package vault

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
)

type (
	// Renewer renews secrets
	Renewer = api.Renewer
)

// Query used by the client to fetch secrets
type Query struct {
	Path    string
	Type    string
	Options map[string]string
}

// Client is a wrapper of Vault API Client
type Client struct {
	client *api.Client
}

// NewClient creates a new client
func NewClient() (*Client, error) {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

// NewRenewer creates a new renewer
func (c *Client) NewRenewer(secret *api.Secret) (*api.Renewer, error) {
	return c.client.NewRenewer(&api.RenewerInput{
		Secret: secret,
	})
}

// GetMount returns a mount
func (c *Client) GetMount(key string) (*api.MountOutput, error) {
	request := c.client.NewRequest("GET", fmt.Sprintf("/v1/sys/internal/ui/mounts/%s", key))

	response, err := c.client.RawRequest(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	props := map[string]interface{}{}
	if err = json.NewDecoder(response.Body).Decode(&props); err != nil {
		return nil, err
	}

	type mount struct {
		Renewable bool            `json:"renewable"`
		Data      api.MountOutput `json:"data" mapstructure:"data"`
	}

	output := &mount{}
	if err = mapstructure.Decode(props, output); err != nil {
		return nil, err
	}

	return &output.Data, nil
}

// List lists path's keys
func (c *Client) List(query *Query) (*api.Secret, error) {
	key := query.Path

	switch query.Type {
	case "kv":
		if version, ok := query.Options["version"]; ok && version == "2" {
			key = path.Join(key, "metadata")
		}
	case "postgresql":
		key = path.Join(key, "roles")
	default:
		return nil, unsupportedTypeError(query.Type)
	}

	return c.client.Logical().List(key)
}

// Read reads path's secrets
func (c *Client) Read(query *Query) (*api.Secret, error) {
	key, err := abs(query.Path, query.Type, query.Options)
	if err != nil {
		return nil, err
	}

	secret, err := c.client.Logical().Read(key)
	if err != nil {
		return nil, err
	}

	switch query.Type {
	case "kv":
		if data, ok := secret.Data["data"].(map[string]interface{}); ok {
			secret.Data = data
		}
	}

	return secret, nil
}

func abs(key, kind string, options map[string]string) (string, error) {
	switch kind {
	case "kv":
		if version, ok := options["version"]; ok && version == "2" {
			dir, name := path.Split(key)
			key = path.Join(dir, "data", name)
		}
	case "postgresql":
		dir, name := path.Split(key)
		if !strings.HasSuffix(dir, "creds/") {
			key = path.Join(dir, "creds", name)
		}
	default:
		return "", unsupportedTypeError(kind)
	}

	return key, nil
}

func unsupportedTypeError(kind string) error {
	return fmt.Errorf("vault: unsupported secret type %s", kind)
}
