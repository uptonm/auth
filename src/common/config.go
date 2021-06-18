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

	Config AppConfig
)

type AppConfig struct {
	Host              string `mapstructure:"HOST"`
	Port              int    `mapstructure:"PORT"`
	Env               Env    `mapstructure:"ENV"`
	Auth0ClientId     string `mapstructure:"AUTH0_CLIENT_ID"`
	Auth0ClientSecret string `mapstructure:"AUTH0_CLIENT_SECRET"`
	Auth0CallbackUrl  string `mapstructure:"AUTH0_CALLBACK_URL"`
	Auth0Domain       string `mapstructure:"AUTH0_DOMAIN"`
	RedisHost         string `mapstructure:"REDIS_HOST"`
	RedisPort         string `mapstructure:"REDIS_PORT"`
	RedisPass         string `mapstructure:"REDIS_PASS"`
}

// Validate accepts a receiver of AppConfig and validates that it includes all required variables
func (c AppConfig) Validate() error {
	if c.Auth0ClientId == "" || c.Auth0ClientSecret == "" || c.Auth0CallbackUrl == "" || c.Auth0Domain == "" {
		return fmt.Errorf("failed to validate config error=auth0 config invalid")
	}

	return nil
}

// ReadConfig utilizes viper to read a common yml file and returns a *AppConfig if its initialized properly
func ReadConfig() error {
	viper.AutomaticEnv()
	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(Config, &envKeysMap); err != nil {
		return err
	}
	for k := range *envKeysMap {
		if err := viper.BindEnv(k); err != nil {
			return err
		}
	}

	viper.SetDefault("host", "localhost")
	viper.SetDefault("port", 8080)
	viper.SetDefault("env", "dev")

	err := viper.Unmarshal(&Config)
	if err != nil {
		return fmt.Errorf("failed to initialize auth0 config")
	}

	err = Config.Validate()
	if err != nil {
		return err
	}

	return nil
}

// IsProd is a helper method accepting *AppConfig as a receiver and resulting in true if within a production environment
func IsProd() bool {
	return Config.Env == Prod
}
