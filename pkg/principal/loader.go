package principal

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/password"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/utils"
	"gopkg.in/yaml.v2"
)

const (
	HtpasswdSeparator  string = ":"
	EmailIndicator     string = "@"
	DefaultEmailDomain string = "localhost"
)

type Loader struct {
	EmailVerified  bool
	EmailDomain    string
	valueProcessor utils.StringProcessor
}

func NewLoader(valueProcessor utils.StringProcessor) *Loader {
	loader := &Loader{
		EmailDomain:    DefaultEmailDomain,
		valueProcessor: valueProcessor,
	}

	return loader
}

func NewPassthroughLoader() *Loader {
	passthrough := func(s string) string {
		return s
	}
	loader := &Loader{
		EmailDomain:    DefaultEmailDomain,
		valueProcessor: utils.FunctorProcessor(passthrough),
	}

	return loader
}

func (l *Loader) ParseMap(data map[string]interface{}) ([]Principal, error) {
	principals := make(map[string]principalData, len(data))
	for k, v := range data {
		d, ok := v.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("malformed map entry %q", k)
		}

		p := &principalData{}
		if err := p.UnmarshalMap(d); err != nil {
			return nil, fmt.Errorf("malformed map entry %q: %w", k, err)
		}

		principals[k] = *p
	}

	return l.parseData(principals)
}

func (l *Loader) ParseFile(path string) ([]Principal, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	principals := make(map[string]principalData)
	err = yaml.Unmarshal(yamlFile, principals)
	if err != nil {
		return nil, fmt.Errorf("unable to parse %q: %w", path, err)
	}

	return l.parseData(principals)
}

func (l *Loader) ParseHtpasswd(path string) ([]Principal, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lineNo := 0
	principals := []Principal{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		principal, err := l.ParseLine(line)
		if err != nil {
			return nil, fmt.Errorf("%s (%d): %w", path, lineNo, err)
		}

		principals = append(principals, principal)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return principals, nil
}

func (l *Loader) ParseLine(line string) (Principal, error) {
	var pass string
	var err error
	principal := &DefaultPrincipal{}
	fields := strings.Split(line, HtpasswdSeparator)

	if len(fields) > 0 && fields[0] != "" {
		principal.Ident = fields[0]
	} else {
		return nil, os.ErrInvalid
	}

	if len(fields) > 1 && fields[1] != "" {
		pass = fields[1]
	} else {
		pass = principal.Ident
	}

	if len(fields) > 2 && fields[2] != "" {
		principal.Email = fields[2]
	} else {
		principal.Email, principal.EmailVerified = l.generateEmail(principal.Ident)
	}

	principal.Password, err = password.ParsePassword(pass)
	if err != nil {
		return nil, err
	}

	return principal, nil
}

func (l *Loader) generateEmail(ident string) (string, bool) {
	if strings.Contains(ident, EmailIndicator) {
		return ident, true
	}

	return fmt.Sprint(ident, EmailIndicator, l.EmailDomain), l.EmailVerified
}

func (l *Loader) parseData(m map[string]principalData) ([]Principal, error) {
	principals := make([]Principal, 0, len(m))

	for i, d := range m {
		principal, err := l.newPrincipal(i, d)
		if err != nil {
			return nil, err
		}

		principals = append(principals, principal)
	}

	return principals, nil
}

func (l *Loader) processPrincipalString(name, value string, field *string) error {
	v, err := l.valueProcessor.Process(value)
	if err != nil {
		return fmt.Errorf("malformed principal %s %q: %w", name, value, err)
	} else if v != "" {
		*field = v
		return nil
	} else if *field == "" {
		return fmt.Errorf("principal %s cannot be empty", name)
	}

	// leave value unchanged
	return nil
}

func (l *Loader) processPrincipalBool(name, value string, field *bool) error {
	v, err := l.valueProcessor.Process(value)
	if err != nil {
		return fmt.Errorf("malformed principal %s %q: %w", name, value, err)
	} else if v == "" {
		// leave value unchanged
		return nil
	}

	b, err := strconv.ParseBool(v)
	if err != nil {
		return fmt.Errorf("malformed principal %s %q: %w", name, v, err)
	}

	*field = b
	return nil
}

func (l *Loader) processPrincipalPassword(name, value string, field *password.Password) error {
	v, err := l.valueProcessor.Process(value)
	if err != nil {
		return fmt.Errorf("malformed principal %s %q: %w", name, value, err)
	} else if v == "" {
		// leave value unchanged
		return nil
	}

	p, err := password.ParsePassword(v)
	if err != nil {
		return fmt.Errorf("malformed principal %s %q: %w", name, v, err)
	}

	field.Value = p.Value
	field.Codec = p.Codec
	return nil
}

func (l *Loader) newPrincipal(ident string, m principalData) (Principal, error) {
	principal := &DefaultPrincipal{}

	if err := l.processPrincipalString("ident", ident, &principal.Ident); err != nil {
		return nil, err
	}

	// calculate default values based on ident
	principal.Password = password.NewPlaintextPassword(principal.Ident)
	principal.Email, principal.EmailVerified = l.generateEmail(principal.Ident)

	if err := l.processPrincipalString("email", m.Email, &principal.Email); err != nil {
		return nil, err
	}

	if err := l.processPrincipalBool("email verified", m.EmailVerified, &principal.EmailVerified); err != nil {
		return nil, err
	}

	if err := l.processPrincipalPassword("password", m.Password, principal.Password); err != nil {
		return nil, err
	}

	return principal, nil
}
