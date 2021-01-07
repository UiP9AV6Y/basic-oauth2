package server

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/square/go-jose.v2"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/log"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/web"
)

func newHealthRouter(config *viper.Viper, logger *log.Controller) (*web.HealthRouter, error) {
	router := web.NewHealthRouter(logger)

	return router, nil
}

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

	routerBuilder := newWebRouter(config)
	routerBuilder.PubKeys = []jose.JSONWebKey{*publicKey}
	routerBuilder.Server = server
	routerBuilder.Signer = signer
	routerBuilder.Logger = logger
	routerBuilder.Login = login
	router, err := routerBuilder.Router()
	if err != nil {
		return nil, err
	}

	return router, nil
}
