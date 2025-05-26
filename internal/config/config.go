package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

var (
	HomeDir    string
	ConfigFile string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get home directory: %v", err)
	}

	HomeDir = home
	ConfigFile = path.Join(home, ".config", "susshi", "config.yaml")
}

type Config struct {
	SSHConfig string `yaml:"ssh_config" mapstructure:"ssh_config"`
	HideIcon  bool   `yaml:"hide_icon" mapstructure:"hide_icon"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigFile(ConfigFile)
	v.SetConfigType("yaml")

	v.SetDefault("ssh_config", path.Join(HomeDir, ".ssh", "config"))

	err := v.ReadInConfig()
	if errors.As(err, &viper.ConfigFileNotFoundError{}) {
		if err := os.MkdirAll(path.Dir(ConfigFile), 0755); err != nil {
			return nil, fmt.Errorf("failed to create config directory: %s", err)
		}

		_, err := os.Create(ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to create config file: %s", err)
		}

		err = v.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %s", err)
		}

	} else if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}

	cfg := &Config{}
	err = v.UnmarshalExact(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %s", err)
	}
	return cfg, nil
}
