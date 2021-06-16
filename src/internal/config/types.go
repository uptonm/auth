package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Host  string      `mapstructure:"host"`
	Port  int         `mapstructure:"port"`
	Auth0 Auth0Config `mapstructure:"auth0"`
}

type Auth0Config struct {
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CallbackUrl  string `mapstructure:"callback_url"`
	Domain       string `mapstructure:"domain"`
}

// Init utilizes viper to read a config yml file and returns a *Config if its initialized properly
func Init() (*Config, error) {
	var config Config

	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
