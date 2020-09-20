package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"

	"github.com/mikeychowy/fiber-crayplate/app/configuration"
	"github.com/mikeychowy/fiber-crayplate/app/providers"
	"github.com/mikeychowy/fiber-crayplate/database"
	"github.com/mikeychowy/fiber-crayplate/routes"
)

func main() {
	// Load configurations
	config, err := configuration.LoadConfigurations()
	if err != nil {
		// Error when loading the configurations
		log.Fatalf("An error occurred while loading the configurations: %v", err)
	}

	// Create a new Fiber application
	app := fiber.New(config.Fiber)

	cb := context.Background()

	// Use the Logger Middleware if enabled
	if config.Enabled["logger"] {
		app.Use(logger.New(config.Logger))
	}

	// Use the Recover Middleware if enabled
	if config.Enabled["recover"] {
		app.Use(recover.New())
	}

	// Use HTTP best practices
	app.Use(func(c *fiber.Ctx) error {
		// Suppress the `www.` at the beginning of URLs
		if config.App.SuppressWWW {
			providers.SuppressWWW(c)
		}
		// Force HTTPS protocol
		if config.App.ForceHTTPS {
			providers.ForceHTTPS(c)
		}
		// Move on the the next route
		return c.Next()
	})

	// Use the Compression Middleware if enabled
	if config.Enabled["compression"] {
		app.Use(compress.New(config.Compression))
	}

	// Use the CORS Middleware if enabled
	if config.Enabled["cors"] {
		app.Use(cors.New(config.CORS))
	}

	// Use the Helmet Middleware if enabled
	if config.Enabled["helmet"] {
		app.Use(helmet.New(config.Helmet))
	}

	// Set hashing provider
	if config.Enabled["hash"] {
		providers.SetHashProvider(config.Hash)
	}

	// Connect to a database
	if config.Enabled["database"] {
		err := database.Connect(*&cb, &config.Database)
		if err != nil {
			exit(&config, app, err)
		}
	}

	// Register application API routes (using the /api/v1 group)
	api := app.Group("/api")
	apiv1 := api.Group("/v1")
	routes.RegisterAPI(apiv1)

	// Set configuration provider
	providers.SetConfiguration(&config)

	// Close any connections on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		exit(&config, app, nil)
	}()

	// Start listening on the specified address
	err = app.Listen(config.App.Listen)
	if err != nil {
		// Exit the application
		exit(&config, app, err)
	}
}

func exit(config *configuration.Configuration, app *fiber.App, err error) {
	// Close database connection
	if config.Enabled["database"] {
		database.Close()
		fmt.Println("Closed database pool")
	}
	// Shutdown Fiber application
	var appErr error
	if err != nil {
		fmt.Printf("Shutting Down Fiber application: %v\n", err)
		appErr = err
	} else {
		appErr = app.Shutdown()
		if appErr != nil {
			fmt.Printf("Fiber application Shutdown Error: %v\n", appErr)
		} else {
			fmt.Println("Fiber application Shutdown.")
		}
	}
	// Return with corresponding exit code
	if appErr != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
