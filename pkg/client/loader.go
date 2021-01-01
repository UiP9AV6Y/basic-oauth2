package client

import (
	"fmt"
	"io/ioutil"

	"github.com/openshift/osin"
	"gopkg.in/yaml.v3"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/utils"
)

type Loader struct {
	valueProcessor utils.StringProcessor
}

func NewLoader(valueProcessor utils.StringProcessor) *Loader {
	loader := &Loader{
		valueProcessor: valueProcessor,
	}

	return loader
}

func NewPassthroughLoader() *Loader {
	passthrough := func(s string) string {
		return s
	}
	loader := &Loader{
		valueProcessor: utils.FunctorProcessor(passthrough),
	}

	return loader
}

func (l *Loader) ParseMap(m map[string]interface{}) ([]osin.Client, error) {
	clients := make(map[string]clientData, len(m))
	for k, v := range m {
		d, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed map entry %q", k)
		}

		c := &clientData{}
		if err := c.UnmarshalMap(d); err != nil {
			return nil, fmt.Errorf("malformed map entry %q: %w", k, err)
		}

		clients[k] = *c
	}

	return l.parseData(clients)
}

func (l *Loader) ParseFile(path string) ([]osin.Client, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	clients := make(map[string]clientData)
	err = yaml.Unmarshal(yamlFile, clients)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %q: %w", path, err)
	}

	return l.parseData(clients)
}

func (l *Loader) parseData(m map[string]clientData) ([]osin.Client, error) {
	clients := make([]osin.Client, 0, len(m))

	for i, d := range m {
		client, err := l.newClient(i, d)
		if err != nil {
			return nil, err
		}

		clients = append(clients, client)
	}

	return clients, nil
}

func (l *Loader) processClientField(name, value string, field *string) error {
	v, err := l.valueProcessor.Process(value)
	if err != nil {
		return fmt.Errorf("malformed client %s %q: %w", name, value, err)
	} else if v == "" {
		return fmt.Errorf("client %s %q cannot be empty", name, value)
	}

	*field = v
	return nil
}

func (l *Loader) newClient(ident string, d clientData) (osin.Client, error) {
	client := &osin.DefaultClient{}

	if err := l.processClientField("ident", ident, &client.Id); err != nil {
		return nil, err
	}

	if err := l.processClientField("secret", d.Secret, &client.Secret); err != nil {
		return nil, err
	}

	if err := l.processClientField("redirect URI", d.RedirectUri, &client.RedirectUri); err != nil {
		return nil, err
	}

	return client, nil
}
