package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// TenantMiddleware ensures team_id is present in context
// This runs after AuthMiddleware and validates the user belongs to a team
func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		teamID := GetTeamID(c)
		if teamID == 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "You must belong to a team to access this resource",
			})
		}
		return c.Next()
	}
}
