package www

import (
	"github.com/gofiber/fiber/v2"
	v1 "github.com/uptonm/auth/src/pkg/v1"
)

// wireRoutes handles the wiring of all packages
func wireRoutes(r fiber.Router) {
	apiGroup := r.Group("/api")
	v1.RegisterV1(apiGroup)
}
