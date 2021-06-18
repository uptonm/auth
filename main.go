package main

import (
	"embed"
	log "github.com/sirupsen/logrus"
	"github.com/uptonm/auth/src/pkg/db"
	"net/http"

	"github.com/uptonm/auth/src/common"
	"github.com/uptonm/auth/src/www"
)

var (
	//go:embed static/template/*
	templates embed.FS
	//go:embed static/res/*
	resources embed.FS

	templateFs http.FileSystem
	resourceFs http.FileSystem
)

func init() {
	var err error
	err = common.ReadConfig()
	if err != nil {
		log.Fatalf("failed to initialize config error=%s", err.Error())
	}

	err = db.InitRedis()
	if err != nil {
		log.Fatalf("failed to connect to redis error=%s", err.Error())
	}
	log.Infof("redis connected")
}

func main() {
	templateFs = common.PickFS(!common.IsProd(), templates, "./static/template")
	resourceFs = common.PickFS(!common.IsProd(), resources, "./static/res")

	www.Start(templateFs, resourceFs)
}
