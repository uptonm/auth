package common

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

type Env string

var (
	Dev  Env = "dev"
	Prod Env = "prod"
)

type Config struct {
	Host              string `mapstructure:"HOST"`
	Port              int    `mapstructure:"PORT"`
	Env               Env    `mapstructure:"ENV"`
	Auth0ClientId     string `mapstructure:"AUTH0_CLIENT_ID"`
	Auth0ClientSecret string `mapstructure:"AUTH0_CLIENT_SECRET"`
	Auth0CallbackUrl  string `mapstructure:"AUTH0_CALLBACK_URL"`
	Auth0Domain       string `mapstructure:"AUTH0_DOMAIN"`
	SigningKey        string `mapstructure:"SIGNING_KEY"`
}

func (c Config) Validate() error {
	if c.Auth0ClientId == "" || c.Auth0ClientSecret == "" || c.Auth0CallbackUrl == "" || c.Auth0Domain == "" {
		return fmt.Errorf("failed to validate config error=auth0 config invalid")
	}

	return nil
}

// ReadConfig utilizes viper to read a common yml file and returns a *Config if its initialized properly
func ReadConfig() (*Config, error) {
	var config Config

	viper.AutomaticEnv()
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(config, &envKeysMap); err != nil {
		return nil, err
	}
	for k := range *envKeysMap {
		if err := viper.BindEnv(k); err != nil {
			return nil, err
		}
	}

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 8080)
	viper.SetDefault("env", "dev")

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize auth0 config")
	}

	err = config.Validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// IsProd is a helper method accepting *Config as a receiver and resulting in true if within a production environment
func (c *Config) IsProd() bool {
	return c.Env == Prod
}
