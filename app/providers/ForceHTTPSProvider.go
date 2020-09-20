package providers

import "github.com/gofiber/fiber/v2"

// ForceHTTPS Forces HTTPS protocol if not forwarded using a reverse proxy
func ForceHTTPS(c *fiber.Ctx) error {
	if c.Get("X-Forwarded-Proto") != "https" && c.Protocol() == "http" {
		return c.Redirect("https://"+c.Hostname()+c.OriginalURL(), 308)
	}
	return nil
}
