package command

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func initRootVersion(config *viper.Viper, flags *flag.FlagSet) {
	flags.BoolP("version", "v", false, "show the application version")
}
