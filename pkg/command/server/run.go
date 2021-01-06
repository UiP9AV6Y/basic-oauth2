package server

import (
	"net/http"
	"os"

	gorilla "github.com/gorilla/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

			logger.Info().Println("Using config file", config.ConfigFileUsed())
			logger.Info().Println("Listening for requests on", server.ListenAddr())
			return server.Run(handler)
		},
		SuggestFor:   []string{"launch", "serve", "start"},
		SilenceUsage: true,
	}

	config.BindEnv("key_file")
	initRunKeySeed(config, cmd.Flags())
	initRunKeySize(config, cmd.Flags())
	initRunKeyGenerate(config, cmd.Flags())
	initRunLogLevel(config, cmd.Flags())
	initRunAllowedAccess(config, cmd.Flags())
	initRunAllowedAuthorize(config, cmd.Flags())
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

	controller := log.NewPlainController(level, os.Stderr)

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

func newHandler(config *viper.Viper, logger *log.Controller) (http.Handler, error) {
	oidc, err := newOIDCRouter(config, logger)
	if err != nil {
		return nil, err
	}

	oidcHandler := oidc.Handler()
	recoverHandler := gorilla.RecoveryHandler(
		gorilla.RecoveryLogger(logger.Fatal()),
	)(oidcHandler)
	logHandler := gorilla.LoggingHandler(os.Stdout, recoverHandler)

	return logHandler, nil
}
