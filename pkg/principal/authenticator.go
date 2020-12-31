package principal

import (
	"errors"
)

var (
	ErrAuthenticationNotFound    = errors.New("no principal found")
	ErrAuthenticationNotPossible = errors.New("principal can not be authenticated")
	ErrAuthenticationNoMatch     = errors.New("principal information does not match input")
)

type Authenticator interface {
	Authenticate(ident, password string) (Principal, error)
}

type Principals map[string]Principal

func (p Principals) AddPrincipal(principal Principal) bool {
	_, ok := p[principal.GetIdent()]
	if ok {
		return false
	}

	p[principal.GetIdent()] = principal
	return true
}

func (p Principals) Authenticate(ident, password string) (Principal, error) {
	principal, ok := p[ident]
	if !ok {
		return nil, ErrAuthenticationNotFound
	}

	pass := principal.GetPassword()
	if pass == nil {
		return principal, ErrAuthenticationNotPossible
	}

	if !pass.Compare(password) {
		return nil, ErrAuthenticationNoMatch
	}

	return principal, nil
}
