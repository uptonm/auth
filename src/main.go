package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/uptonm/auth/src/internal/config"
	"github.com/uptonm/auth/src/www"
)

var Config *config.Config

func init() {
	var err error
	Config, err = config.Init()
	if err != nil {
		log.Fatalf("failed to initialize config error=%s", err.Error())
	}
}

func main() {
	www.Start(Config)
}
