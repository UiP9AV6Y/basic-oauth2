package principal

import (
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

var (
	ErrUnmarshalPrincipal     = errors.New("unable to unmarshal principal")
	ErrUnmarshalPassword      = errors.New("unable to unmarshal password")
	ErrUnmarshalEmail         = errors.New("unable to unmarshal email")
	ErrUnmarshalEmailVerified = errors.New("email_verified is neither a string nor a boolean value")
)

type principalData struct {
	Password      string `yaml:"password,omitempty"`
	Email         string `yaml:"email,omitempty"`
	EmailVerified string `yaml:"email_verified,omitempty"`
}

func (d *principalData) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.MappingNode {
		data := make(map[string]interface{})
		err := value.Decode(data)
		if err != nil {
			return err
		}

		return d.UnmarshalMap(data)
	} else if value.Kind == yaml.ScalarNode {
		return value.Decode(&d.Password)
	}

	return ErrUnmarshalPrincipal
}

func (d *principalData) UnmarshalMap(data map[string]interface{}) error {
	for k, v := range data {
		switch k {
		case "password":
			if s, ok := v.(string); ok {
				d.Password = s
			} else {
				return ErrUnmarshalPassword
			}
		case "email":
			if s, ok := v.(string); ok {
				d.Email = s
			} else {
				return ErrUnmarshalEmail
			}
		case "email_verified":
			if s, ok := v.(string); ok {
				d.EmailVerified = s
			} else if b, ok := v.(bool); ok {
				d.EmailVerified = strconv.FormatBool(b)
			} else {
				return ErrUnmarshalEmailVerified
			}
		default:
			return fmt.Errorf("%q is not a valid attribute", k)
		}
	}

	return nil
}
