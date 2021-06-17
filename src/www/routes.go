package www

import (
	"github.com/uptonm/auth/src/pkg/auth"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	"github.com/uptonm/auth/src/common"
	v1 "github.com/uptonm/auth/src/pkg/v1"
)

// wireRoutes handles the wiring of all packages
func wireRoutes(r fiber.Router, resourceFs http.FileSystem, config *common.Config) {
	apiGroup := r.Group("/api")
	auth.RegisterAuthRoutes(r)
	v1.RegisterV1(apiGroup, config)

	// Predefined route for favicon at root of domain
	r.Use(favicon.New(favicon.Config{
		File:       "favicon.ico",
		FileSystem: resourceFs,
	}))

	// Serve static files from /static/resources preventing directory listings
	r.Use(filesystem.New(filesystem.Config{
		Root:   common.StrictFs{Fs: resourceFs},
		MaxAge: 86400,
	}))

	NotFound(r)
}

func NotFound(r fiber.Router) {
	r.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).Render("404", map[string]string{
			"PageName": "Page Not Found",
		})
	})
}
