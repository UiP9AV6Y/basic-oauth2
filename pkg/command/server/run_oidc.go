package server

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/square/go-jose.v2"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/log"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/web"
)

func newOIDCRouter(config *viper.Viper, logger *log.Controller) (*web.OIDCRouter, error) {
	server, err := newOauth2Server(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 server: %w", err)
	}

	keyname, secret, err := newOauth2Secret(config, logger)
	if err != nil {
		return nil, fmt.Errorf("secret %q is not parseable: %w", keyname, err)
	}

	signer, err := secret.NewSigner()
	if err != nil {
		return nil, fmt.Errorf("failed to create jwtSigner: %w", err)
	}

	publicKey, err := secret.NewPublicKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}

	login, err := newAuthenticator(config)
	if err != nil {
		return nil, err
	}

	handlerBuilder := newWebRouter(config)
	handlerBuilder.PubKeys = []jose.JSONWebKey{*publicKey}
	handlerBuilder.Server = server
	handlerBuilder.Signer = signer
	handlerBuilder.Logger = logger
	handlerBuilder.Login = login
	handler, err := handlerBuilder.Router()
	if err != nil {
		return nil, err
	}

	return handler, nil
}
