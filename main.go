package main

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/uptonm/auth/src/common"
	"github.com/uptonm/auth/src/www"
)

var (
	Config *common.Config

	//go:embed static/template/*
	templates embed.FS
	//go:embed static/res/*
	resources embed.FS

	templateFs http.FileSystem
	resourceFs http.FileSystem
)

func init() {
	var err error
	Config, err = common.ReadConfig()
	if err != nil {
		log.Fatalf("failed to initialize config error=%s", err.Error())
	}
}

func main() {
	templateFs = common.PickFS(!Config.IsProd(), templates, "./static/template")
	resourceFs = common.PickFS(!Config.IsProd(), resources, "./static/res")

	www.Start(Config, templateFs, resourceFs)
}
