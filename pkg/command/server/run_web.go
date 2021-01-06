package server

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/web"
)

func newWebServer(config *viper.Viper) *web.ServerOptions {
	server := &web.ServerOptions{
		Socket:  config.GetString("listen_socket"),
		Address: config.GetString("listen_address"),
		Port:    config.GetInt("listen_port"),
	}

	return server
}

func newWebRouter(config *viper.Viper) *web.OIDCRouterOptions {
	address := config.GetString("listen_address")
	if address == "0.0.0.0" {
		address = "127.0.0.1"
	} else if address == "::" {
		address = "::1"
	}

	port := config.GetInt("listen_port")
	if port <= 0 {
		port = 80
	}

	defaultUrl := fmt.Sprintf("http://%s:%d/", address, port)
	issuerUrl := config.GetString("issuer_url")
	if issuerUrl == "" {
		issuerUrl = defaultUrl
	}

	authBaseUrl := config.GetString("auth_base_url")
	if authBaseUrl == "" {
		authBaseUrl = issuerUrl
	}

	infoBaseUrl := config.GetString("info_base_url")
	if infoBaseUrl == "" {
		infoBaseUrl = issuerUrl
	}

	tokenBaseUrl := config.GetString("token_base_url")
	if tokenBaseUrl == "" {
		tokenBaseUrl = issuerUrl
	}

	keysBaseUrl := config.GetString("keys_base_url")
	if keysBaseUrl == "" {
		keysBaseUrl = issuerUrl
	}

	handler := &web.OIDCRouterOptions{
		IssuerUrl:    issuerUrl,
		AuthBaseUrl:  authBaseUrl,
		InfoBaseUrl:  infoBaseUrl,
		TokenBaseUrl: tokenBaseUrl,
		KeysBaseUrl:  keysBaseUrl,
	}

	return handler
}
