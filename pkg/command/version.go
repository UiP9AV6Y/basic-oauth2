package command

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/UiP9AV6Y/basic-oauth2/pkg/version"
)

func NewVersionCommand(name string, _ *viper.Viper) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the application version",
		Long:  `Print the application version to stdout and exit.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s version %s\n", name, version.Version())
			fmt.Printf("commit: %s\n", version.Commit())

			return nil
		},
		SuggestFor:   []string{"info"},
		SilenceUsage: true,
	}

	return cmd
}
