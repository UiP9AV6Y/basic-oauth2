package server

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewHealthCommand(config *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Server healthcheck",
		Long:  `Test the server healthyness.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		SuggestFor:   []string{"launch", "serve", "start"},
		SilenceUsage: true,
	}

	return cmd
}
