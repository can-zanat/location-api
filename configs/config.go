package configs

import (
	"github.com/spf13/viper"
)

type Config struct {
	MongoDB struct {
		URI string `mapstructure:"uri"`
	} `mapstructure:"mongoDB"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
