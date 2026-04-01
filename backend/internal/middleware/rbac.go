package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// Admin roles that can manage team resources
var adminRoles = map[string]bool{
	"captain":      true,
	"vice_captain": true,
	"manager":      true,
	"strategist":   true,
}

// RequireAdmin restricts access to admin roles only
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := GetRole(c)
		if !adminRoles[role] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}
		return c.Next()
	}
}

// RequireRole restricts access to specific roles
func RequireRole(roles ...string) fiber.Handler {
	allowed := make(map[string]bool)
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *fiber.Ctx) error {
		role := GetRole(c)
		if !allowed[role] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}
		return c.Next()
	}
}

// IsAdmin checks if the current user has an admin role (helper, not middleware)
func IsAdmin(role string) bool {
	return adminRoles[role]
}
