package command

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func initRootVerbose(config *viper.Viper, flags *flag.FlagSet) {
	flags.BoolP("verbose", "V", false, "verbose output")
	config.BindPFlag("verbose", flags.Lookup("verbose"))
	config.BindEnv("verbose")
}

func initRootVersion(config *viper.Viper, flags *flag.FlagSet) {
	flags.BoolP("version", "v", false, "show the application version")
}
