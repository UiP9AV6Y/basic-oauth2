package client

import (
	"errors"
	"fmt"
)

var (
	ErrUnmarshalSecret      = errors.New("unable to unmarshal secret")
	ErrUnmarshalRedirectURI = errors.New("unable to unmarshal redirect_uri")
)

type clientData struct {
	Secret      string `yaml:"secret,omitempty"`
	RedirectUri string `yaml:"redirect_uri,omitempty"`
}

func (d *clientData) UnmarshalMap(data map[string]interface{}) error {
	for k, v := range data {
		switch k {
		case "secret":
			if s, ok := v.(string); ok {
				d.Secret = s
			} else {
				return ErrUnmarshalSecret
			}
		case "redirect_uri":
			if s, ok := v.(string); ok {
				d.RedirectUri = s
			} else {
				return ErrUnmarshalRedirectURI
			}
		default:
			return fmt.Errorf("%q is not a valid attribute", k)
		}
	}

	return nil
}
