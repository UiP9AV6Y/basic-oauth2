package server

import (
	"fmt"
	l "log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/square/go-jose.v2"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/log"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/utils"
	"github.com/UiP9AV6Y/basic-oauth2/pkg/web"
)

func NewRunCommand(config *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run KEYFILE",
		Short: "Run the server",
		Long:  `Run the application server`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				config.Set("key_file", args[0])
			}

			logger, err := newLogger(config)
			if err != nil {
				return err
			}

			server, err := newServer(config)
			if err != nil {
				return err
			}

			handler, err := newHandler(config, logger)
			if err != nil {
				return err
			}

			logger.Info().Printf("Listening for requests on %q", server.ListenAddr())
			return server.Run(handler.Handler())
		},
		SuggestFor:   []string{"launch", "serve", "start"},
		SilenceUsage: true,
	}

	config.BindEnv("key_file")
	initRunKeySeed(config, cmd.Flags())
	initRunKeySize(config, cmd.Flags())
	initRunKeyGenerate(config, cmd.Flags())
	initRunLogLevel(config, cmd.Flags())
	initRunAuthorizationExpiration(config, cmd.Flags())
	initRunAccessExpiration(config, cmd.Flags())
	initRunErrorStatusCode(config, cmd.Flags())
	initRunAllowClientSecretInParams(config, cmd.Flags())
	initRunAllowGetAccessRequest(config, cmd.Flags())
	initRunRequirePKCEForPublicClients(config, cmd.Flags())
	initRunRedirectUriSeparator(config, cmd.Flags())
	initRunRetainTokenAfterRefresh(config, cmd.Flags())
	initRunIssuerUrl(config, cmd.Flags())
	initRunAuthBaseUrl(config, cmd.Flags())
	initRunInfoBaseUrl(config, cmd.Flags())
	initRunTokenBaseUrl(config, cmd.Flags())
	initRunKeysBaseUrl(config, cmd.Flags())
	initRunEmailDomain(config, cmd.Flags())

	return cmd
}

func newStringProcessor(config *viper.Viper) utils.StringProcessor {
	return utils.NewChainProcessor(
		utils.NewExpandingProcessor(utils.EnvProcessor),
		utils.NewRelativeFileReferenceProcessor('!', config.ConfigFileUsed()),
	)
}

func newLogger(config *viper.Viper) (*log.Controller, error) {
	name := config.GetString("log_level")
	level, err := log.ParseLevelName(name)
	if err != nil {
		return nil, err
	}

	logger := l.New(os.Stderr, "", 0)
	controller := log.NewController(level, logger)

	return controller, nil
}

func newServer(config *viper.Viper) (*web.Server, error) {
	serverBuilder := newWebServer(config)
	server, err := serverBuilder.Server()
	if err != nil {
		return nil, err
	}

	return server, nil
}

func newHandler(config *viper.Viper, logger *log.Controller) (*web.HttpHandler, error) {
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

	handlerBuilder := newWebHandler(config)
	handlerBuilder.PubKeys = []jose.JSONWebKey{*publicKey}
	handlerBuilder.Server = server
	handlerBuilder.Signer = signer
	handlerBuilder.Logger = logger
	handlerBuilder.Login = login
	handler, err := handlerBuilder.Handler()
	if err != nil {
		return nil, err
	}

	return handler, nil
}
