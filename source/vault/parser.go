package vault

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/phogolabs/cli"
)

var _ cli.Parser = &Parser{}

// Parser is a parser that populates flags from Hashi Corp Vault
type Parser struct {
	client  *Client
	renwers []*api.Renewer
}

// Parse parses the args
func (m *Parser) Parse(ctx *cli.Context) error {
	client, err := NewClient()
	if err != nil {
		return err
	}

	m.client = client

	for _, flag := range ctx.Command.Flags {
		definition := flag.Definition()

		key, ok := definition.Metadata["vault_key"]
		if !ok {
			continue
		}

		data, err := m.get(key)
		if err != nil {
			return err
		}

		value, err := m.transform(data)
		if err != nil {
			return err
		}

		if err := flag.Set(value); err != nil {
			return err
		}
	}

	return nil
}

func (m *Parser) transform(data interface{}) (string, error) {
	buffer := &bytes.Buffer{}

	if err := json.NewEncoder(buffer).Encode(data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (m *Parser) get(root string) (interface{}, error) {
	mnt, err := m.client.GetMount(root)
	if err != nil {
		return nil, err
	}

	keys, err := m.keys(root, mnt)
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		keys = append(keys, root)
	}

	var data []map[string]interface{}

	for _, key := range keys {
		secret, err := m.fetch(key, mnt)
		if err != nil {
			return nil, err
		}

		data = append(data, secret.Data)

		if key == root {
			return data[0], nil
		}
	}

	return data, nil
}

func (m *Parser) fetch(key string, mnt *api.MountOutput) (*api.Secret, error) {
	secret, err := m.client.Read(&Query{
		Path:    key,
		Type:    mnt.Type,
		Options: mnt.Options,
	})

	if err != nil {
		return nil, err
	}

	renewer, err := m.client.NewRenewer(secret)
	if renewer != nil {
		m.renwers = append(m.renwers, renewer)
	}

	return secret, nil
}

func (m *Parser) keys(key string, opt *api.MountOutput) ([]string, error) {
	keys := []string{}

	list, err := m.client.List(&Query{
		Path:    key,
		Type:    opt.Type,
		Options: opt.Options,
	})

	if list == nil || err != nil {
		return keys, err
	}

	data, ok := list.Data["keys"].([]interface{})
	if !ok {
		return keys, fmt.Errorf("vault: keys are missing for mount %s", key)
	}

	for _, secret := range data {
		path := fmt.Sprintf("%s/%v", key, secret)
		keys = append(keys, path)
	}

	return keys, nil
}
