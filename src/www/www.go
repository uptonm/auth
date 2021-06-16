package www

import (
	"fmt"
	"github.com/gofiber/template/html"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/auth0"
	log "github.com/sirupsen/logrus"

	"github.com/uptonm/auth/src/internal/config"
)

// Start handles all initialization of the go-fiber webserver
func Start(config *config.Config) {
	engine := html.New("./static", ".html")

	r := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})

	goth.UseProviders(
		auth0.New(config.Auth0.ClientId, config.Auth0.ClientSecret, config.Auth0.CallbackUrl, config.Auth0.Domain),
	)

	r.Get("/", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{})
	})

	wireMiddleware(r)
	wireRoutes(r)

	// graceful shutdown with SIGINT | SIGTERM and others will hard kill
	// credit for this lovely method https://github.com/dechristopher/dchr.host
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		_ = r.Shutdown()
	}()

	// listen for connections on primary listening port
	log.Infof("auth listening on %s:%d", config.Host, config.Port)
	if err := r.Listen(fmt.Sprintf("%s:%d", config.Host, config.Port)); err != nil {
		log.Error(err.Error())
	}

	// exit cleanly
	os.Exit(0)
}
