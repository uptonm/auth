package v1

import (
	"github.com/gofiber/fiber/v2"
	auth2 "github.com/uptonm/auth/src/pkg/auth"
	"github.com/uptonm/auth/src/pkg/v1/health"
	"github.com/uptonm/auth/src/pkg/v1/protected"
)

// RegisterV1 handles the registration of all routes of the v1 package
func RegisterV1(r fiber.Router) {
	v1Group := r.Group("/v1")
	health.RegisterHealth(v1Group)
	auth2.RegisterAuthRoutes(v1Group)
	protected.RegisterProtected(v1Group)
}
