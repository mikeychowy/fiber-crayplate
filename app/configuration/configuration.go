package configuration

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"

	hashing "github.com/thomasvvugt/fiber-hashing"
)

// Configuration struct of each config type
type Configuration struct {
	Fiber          fiber.Config
	App            ApplicationConfiguration
	Enabled        map[string]bool
	Logger         logger.Config
	TemplateEngine func(raw string, bind interface{}) (out string, err error)
	Compression    compress.Config
	CORS           cors.Config
	Helmet         helmet.Config
	Hash           hashing.Config
	PublicPrefix   string
	PublicRoot     string
	Public         fiber.Static
	Database       DatabaseConfiguration
}

// LoadConfigurations using viper
func LoadConfigurations() (config Configuration, err error) {
	config.Enabled = make(map[string]bool)
	// Load the Fiber application configuration
	fiberSettings, err := loadFiberConfiguration()
	if err != nil {
		return config, err
	}
	config.Fiber = fiberSettings

	// Load the application configuration
	appConfig, err := loadApplicationConfiguration()
	if err != nil {
		return config, err
	}
	config.App = appConfig

	// Load the logger middleware configuration
	loggerEnabled, loggerConfig, err := loadLoggerConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["logger"] = loggerEnabled
	config.Logger = loggerConfig

	// Load the recover middleware configuration
	recoverEnabled, err := loadRecoverConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["recover"] = recoverEnabled

	// Load the compression middleware configuration
	compressionEnabled, compressionConfig, err := loadCompressionConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["compression"] = compressionEnabled
	config.Compression = compressionConfig

	// Load the CORS middleware configuration
	corsEnabled, corsConfig, err := loadCORSConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["cors"] = corsEnabled
	config.CORS = corsConfig

	// Load the Helmet middleware configuration
	helmetEnabled, helmetConfig, err := loadHelmetConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["helmet"] = helmetEnabled
	config.Helmet = helmetConfig

	// Load the hashing configuration
	hashEnabled, hashConfig, err := loadHashConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["hash"] = hashEnabled
	config.Hash = hashConfig

	// Load the database configuration
	databaseEnabled, databaseConfig, err := loadDatabaseConfiguration()
	if err != nil {
		return config, err
	}
	config.Enabled["database"] = databaseEnabled
	config.Database = databaseConfig

	// Return the configuration
	return config, nil
}
