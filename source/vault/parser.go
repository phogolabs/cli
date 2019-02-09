package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"
	"github.com/phogolabs/cli"
)

var _ cli.Parser = &Parser{}

// Parser is a parser that populates flags from Hashi Corp Vault
type Parser struct {
	client *api.Client
}

// Parse parses the args
func (m *Parser) Parse(ctx *cli.Context) error {
	if err := m.init(ctx); err != nil {
		return err
	}

	if m.client == nil {
		return nil
	}

	for _, flag := range ctx.Command.Flags {
		key := key(flag.Definition().Metadata)

		if key == "" {
			continue
		}

		mnt, err := m.mount(key)
		if err != nil {
			return err
		}

		key = abs(key, mnt)

		secret, err := m.client.Logical().Read(key)
		if err != nil {
			return err
		}

		if err = m.rewnew(secret); err != nil {
			return err
		}

		value, err := transform(secret.Data, mnt)
		if err != nil {
			return err
		}

		if err := flag.Set(value); err != nil {
			return err
		}
	}

	return nil
}

func (m *Parser) init(ctx *cli.Context) error {
	config := api.DefaultConfig()

	if addr := ctx.String("vault-addr"); addr != "" {
		config.Address = addr
	}

	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	if token := ctx.String("vault-token"); token != "" {
		m.client = client
		client.SetToken(token)
		return nil
	}

	//TODO: kubernetes

	return nil
}

func (m *Parser) mount(key string) (*api.MountOutput, error) {
	path := fmt.Sprintf("/v1/sys/internal/ui/mounts/%s", key)
	request := m.client.NewRequest("GET", path)

	response, err := m.client.RawRequest(request)
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

func (m *Parser) rewnew(secret *api.Secret) error {
	renewer, err := m.client.NewRenewer(&api.RenewerInput{
		Secret: secret,
	})

	if err != nil {
		return err
	}

	go renewer.Renew()
	return nil
}

func key(metadata map[string]string) string {
	if metadata == nil {
		return ""
	}

	if key, ok := metadata["vault_key"]; ok {
		return key
	}

	return ""
}

func abs(key string, mnt *api.MountOutput) string {
	switch mnt.Type {
	case "kv":
		if version, ok := mnt.Options["version"]; ok && version == "2" {
			dir, name := path.Split(key)
			key = path.Join(dir, "data", name)
		}
	}

	return key
}

func transform(data map[string]interface{}, mnt *api.MountOutput) (string, error) {
	buffer := &bytes.Buffer{}

	var value interface{}

	switch mnt.Type {
	case "kv":
		if version, ok := mnt.Options["version"]; ok && version == "2" {
			value = data["data"]
		}
	default:
		value = data
	}

	if err := json.NewEncoder(buffer).Encode(value); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
