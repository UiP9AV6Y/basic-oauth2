package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/web"
)

func NewHealthCommand(config *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Server healthcheck",
		Long:  `Test the server healthyness.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(config)
			if err != nil {
				return err
			}

			consumer, err := newConsumer(config)
			if err != nil {
				return err
			}

			return client.Visit(web.HealthEndpoint, consumer)
		},
		SuggestFor:   []string{"launch", "serve", "start"},
		SilenceUsage: true,
	}

	return cmd
}

func newConsumer(config *viper.Viper) (web.HealthConsumer, error) {
	consumer := web.NewAPIHealthConsumer()

	return consumer, nil
}

func newClient(config *viper.Viper) (*web.Client, error) {
	clientBuilder := newWebClient(config)
	client, err := clientBuilder.Client()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func newWebClient(config *viper.Viper) *web.ClientOptions {
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

	server := &web.ClientOptions{
		Socket:  config.GetString("listen_socket"),
		Address: address,
		Port:    port,
	}

	return server
}
