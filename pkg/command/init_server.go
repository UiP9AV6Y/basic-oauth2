package command

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ServerPort    int    = 9096
	ServerAddress string = "0.0.0.0"
)

func initServerPort(config *viper.Viper, flags *flag.FlagSet) {
	flags.Int("port", ServerPort, "Port to bind to")
	_ = config.BindPFlag("listen_port", flags.Lookup("port"))
	_ = config.BindEnv("listen_port")
}

func initServerAddress(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("address", ServerAddress, "Address to bind to")
	_ = config.BindPFlag("listen_address", flags.Lookup("address"))
	_ = config.BindEnv("listen_address")
}

func initServerSocket(config *viper.Viper, flags *flag.FlagSet) {
	flags.String("socket", "", "Unix socket to bind to")
	_ = config.BindPFlag("listen_socket", flags.Lookup("socket"))
	_ = config.BindEnv("listen_socket")
}
