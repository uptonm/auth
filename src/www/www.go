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
func Start(templateFs, resourceFs http.FileSystem) {
	// populate template engine from templates filesystem
	engine = html.NewFileSystem(templateFs, ".html")

	// enable template engine reloading on dev
	engine.Reload(!common.IsProd())

	r := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})

	goth.UseProviders(
		auth0.New(common.Config.Auth0ClientId, common.Config.Auth0ClientSecret, common.Config.Auth0CallbackUrl, common.Config.Auth0Domain),
	)

	r.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"PageName": "Portfolio",
		})
	})

	wireMiddleware(r)
	wireRoutes(r, resourceFs)

	// graceful shutdown with SIGINT | SIGTERM and others will hard kill
	// credit for this lovely method https://github.com/dechristopher/dchr.host
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		_ = r.Shutdown()
	}()

	// listen for connections on primary listening port
	log.Infof("uptonm.io listening on %s:%d", common.Config.Host, common.Config.Port)
	if err := r.Listen(fmt.Sprintf("%s:%d", common.Config.Host, common.Config.Port)); err != nil {
		log.Error(err.Error())
	}

	// exit cleanly
	os.Exit(0)
}
