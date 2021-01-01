package command

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewRootCommand(name string, config *viper.Viper) *cobra.Command {
	srvCmd := NewServerCommand(name, config)
	verCmd := NewVersionCommand(name, config)
	cmd := &cobra.Command{
		Use:   name,
		Short: "oAuth2 service with Basic authentication",
		Long: `HTTP server providing OIDC endpoints for oAuth2
authentication flows using Basic authentication
as means of identifying principals.`,
		RunE:         versionFlagProxy(verCmd),
		SilenceUsage: true,
	}

	initRootVersion(config, cmd.Flags())

	cmd.AddCommand(srvCmd)
	cmd.AddCommand(verCmd)

	return cmd
}

func versionFlagProxy(versionCmd *cobra.Command) func(cmd *cobra.Command, args []string) error {
	proxy := func(cmd *cobra.Command, args []string) error {
		versionVal, err := cmd.Flags().GetBool("version")
		if err != nil {
			return err
		} else if versionVal {
			return versionCmd.RunE(cmd, args)
		}

		return flag.ErrHelp
	}

	return proxy
}
