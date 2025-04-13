package configs

import (
	"fmt"

	"github.com/spf13/viper"
)

func NewConfig(vars *Vars) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetConfigFile(fmt.Sprintf("%s.yaml", vars.Profile))

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	return &config, err
}
