package www

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/auth0"
	log "github.com/sirupsen/logrus"

	"github.com/uptonm/auth/src/common"
)

var (
	// fiber html template engine
	engine *html.Engine
)

// Start handles all initialization of the go-fiber webserver
func Start(config *common.Config, templateFs, resourceFs http.FileSystem) {
	// populate template engine from templates filesystem
	engine = html.NewFileSystem(templateFs, ".html")

	// enable template engine reloading on dev
	engine.Reload(!config.IsProd())

	r := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})

	goth.UseProviders(
		auth0.New(config.Auth0ClientId, config.Auth0ClientSecret, config.Auth0CallbackUrl, config.Auth0Domain),
	)

	r.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"PageName": "Portfolio",
		})
	})

	wireMiddleware(r)
	wireRoutes(r, resourceFs, config)

	// graceful shutdown with SIGINT | SIGTERM and others will hard kill
	// credit for this lovely method https://github.com/dechristopher/dchr.host
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		_ = r.Shutdown()
	}()

	// listen for connections on primary listening port
	log.Infof("uptonm.io listening on %s:%d", config.Host, config.Port)
	if err := r.Listen(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
		log.Error(err.Error())
	}

	// exit cleanly
	os.Exit(0)
}
