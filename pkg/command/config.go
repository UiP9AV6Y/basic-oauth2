package command

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	SystemConfigDir string = "/etc"
	ConfigName      string = "config"
	ConfigFormat    string = "yaml"
)

func NewConfig(exec string) (name string, config *viper.Viper, err error) {
	var cfgDir string
	cfgDir, err = os.UserConfigDir()
	if err != nil {
		return
	}

	config = viper.New()
	name = filepath.Base(exec)

	replace := strings.NewReplacer("-", "_", ".", "_")
	prefix := strings.ToUpper(replace.Replace(name))
	config.SetEnvKeyReplacer(replace)
	config.SetEnvPrefix(name)
	config.AllowEmptyEnv(true)

	if cfgFile := os.Getenv(prefix + "_CONFIG_FILE"); cfgFile != "" {
		config.SetConfigFile(cfgFile)
	} else {
		sysDir := filepath.Join(SystemConfigDir, name)
		usrDir := filepath.Join(cfgDir, name)
		exeDir := filepath.Dir(exec)

		config.SetConfigName(ConfigName)
		config.SetConfigType(ConfigFormat)
		config.AddConfigPath(sysDir)
		config.AddConfigPath(usrDir)
		config.AddConfigPath(exeDir)
		config.AddConfigPath(".")
	}

	return
}
