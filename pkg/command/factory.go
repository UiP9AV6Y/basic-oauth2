package command

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func BasicOauth2() error {
	exec, err := os.Executable()
	if err != nil {
		return err
	}

	name, config, err := NewConfig(exec)
	if err != nil {
		return err
	}

	cmd := NewRootCommand(name, config)

	err = config.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("Unable to read config file: %w", err)
		}
	}

	return cmd.Execute()
}
