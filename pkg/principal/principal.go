package principal

import (
	"github.com/UiP9AV6Y/basic-oauth2/pkg/password"
)

type Principal interface {
	// GetIdent returns the principal identifier.
	GetIdent() string
	// GetPassword returns the authentication secret.
	GetPassword() *password.Password
	// GetEmail returns the email address associated with
	// this principal as well as the verification state.
	GetEmail() (string, bool)
}

type DefaultPrincipal struct {
	Password      *password.Password
	Ident         string
	Email         string
	EmailVerified bool
}

func (p *DefaultPrincipal) GetIdent() string {
	return p.Ident
}

func (p *DefaultPrincipal) GetPassword() *password.Password {
	return p.Password
}

func (p *DefaultPrincipal) GetEmail() (string, bool) {
	return p.Email, p.EmailVerified
}

func NewEmailPrincipal(email string, pass *password.Password) *DefaultPrincipal {
	principal := &DefaultPrincipal{
		Password:      pass,
		Ident:         email,
		Email:         email,
		EmailVerified: false,
	}
	return principal
}
