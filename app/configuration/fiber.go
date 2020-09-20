package configuration

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/spf13/viper"
)

func loadFiberConfiguration() (settings fiber.Config, err error) {
	// Set a new configuration provider
	provider := viper.New()

	// Set configuration provider settings
	provider.SetConfigName("fiber")
	provider.AddConfigPath("./config")

	// Set default configurations
	setDefaultFiberConfiguration(provider)

	// Read configuration file
	err = provider.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error since we have default configurations
		} else {
			// Config file was found but another error was produced
			return settings, err
		}
	}

	// Unmarshal the configuration file into fiber.Settings
	err = provider.Unmarshal(&settings)

	// Return the configuration (and error if occurred)
	return settings, err
}

// Set default configuration for Fiber
func setDefaultFiberConfiguration(provider *viper.Viper) {
	provider.SetDefault("Prefork", false)
	provider.SetDefault("ServerHeader", "")
	provider.SetDefault("StrictRouting", false)
	provider.SetDefault("CaseSensitive", false)
	provider.SetDefault("Immutable", false)
	provider.SetDefault("UnescapePath", false)
	provider.SetDefault("BodyLimit", 4*1024*1024)
	provider.SetDefault("Concurrency", 256*1024)
	provider.SetDefault("DisableKeepalive", false)
	provider.SetDefault("DisableDefaultDate", false)
	provider.SetDefault("DisableDefaultContentType", false)
	provider.SetDefault("DisableStartupMessage", false)
	provider.SetDefault("ETag", false)
	provider.SetDefault("ReadTimeout", nil)
	provider.SetDefault("WriteTimeout", nil)
	provider.SetDefault("IdleTimeout", nil)
	provider.SetDefault("ReadBufferSize", 4096)
	provider.SetDefault("WriteBufferSize", 4096)
	provider.SetDefault("CompressedFileSuffix", ".fiber.gz")
	provider.SetDefault("ProxyHeader", "")
	provider.SetDefault("GETOnly", false)
	provider.SetDefault("ErrorHandler", defaultAPIErrorHandler)
}

// Default Error Handler
var defaultAPIErrorHandler = func(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	c.Status(code)
	err = c.JSON(fiber.Map{
		"success": false,
		"status":  code,
		"message": fmt.Sprintf("%s", err),
		"data":    make([]int, 0, 1),
	})
	if err != nil {
		return c.Status(500).SendString("Internal Server Error")
	}
	return nil
}
