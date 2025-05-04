package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewConfig(vars *Vars) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigFile(fmt.Sprintf("%s.yaml", vars.Profile))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
