package server

import (
	"github.com/openshift/osin"
)

type authorizeTypes map[string]osin.AuthorizeRequestType

func (t authorizeTypes) FlagValues() []string {
	var i int
	values := make([]string, len(t))

	for k, _ := range t {
		values[i] = k
		i++
	}

	return values
}

var allowedAuthorizeTypes = authorizeTypes{
	"code":  osin.CODE,
	"token": osin.TOKEN,
}

type accessTypes map[string]osin.AccessRequestType

func (t accessTypes) FlagValues() []string {
	var i int
	values := make([]string, len(t))

	for k, _ := range t {
		values[i] = k
		i++
	}

	return values
}

var allowedAccessTypes = accessTypes{
	"assertion":          osin.ASSERTION,
	"authorization-code": osin.AUTHORIZATION_CODE,
	"client-credentials": osin.CLIENT_CREDENTIALS,
	"password":           osin.PASSWORD,
	"refresh-token":      osin.REFRESH_TOKEN,
}
