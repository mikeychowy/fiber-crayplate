package providers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// SuppressWWW suppresses the `www.` at the beginning of URLs
func SuppressWWW(c *fiber.Ctx) error {
	hostnameSplit := strings.Split(c.Hostname(), ".")
	if hostnameSplit[0] == "www" && len(hostnameSplit) > 1 {
		newHostname := ""
		for i := 1; i <= (len(hostnameSplit) - 1); i++ {
			if i != (len(hostnameSplit) - 1) {
				newHostname = newHostname + hostnameSplit[i] + "."
			} else {
				newHostname = newHostname + hostnameSplit[i]
			}
		}
		return c.Redirect(c.Protocol()+"://"+newHostname+c.OriginalURL(), 301)
	}
	return nil
}
