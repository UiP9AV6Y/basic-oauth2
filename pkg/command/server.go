package command

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/command/server"
)

func NewServerCommand(name string, config *viper.Viper) *cobra.Command {
	healthCmd := server.NewHealthCommand(config)
	runCmd := server.NewRunCommand(config)
	cmd := &cobra.Command{
		Use:          "server",
		Short:        runCmd.Short,
		Long:         runCmd.Long,
		SuggestFor:   []string{"web", "listener", "daemon"},
		SilenceUsage: true,
	}

	initServerPort(config, cmd.PersistentFlags())
	initServerAddress(config, cmd.PersistentFlags())
	initServerSocket(config, cmd.PersistentFlags())

	cmd.AddCommand(healthCmd)
	cmd.AddCommand(runCmd)

	return cmd
}
