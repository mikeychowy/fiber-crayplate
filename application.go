package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/helmet"

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
		os.Exit(2)
	}

	// Create a new Fiber application
	app := fiber.New(&config.Fiber)

	cb := context.Background()

	// Use the Logger Middleware if enabled
	if config.Enabled["logger"] {
		app.Use(middleware.Logger(config.Logger))
	}

	// Use the Recover Middleware if enabled
	if config.Enabled["recover"] {
		app.Use(middleware.Recover())
	}

	// Use HTTP best practices
	app.Use(func(c *fiber.Ctx) {
		// Suppress the `www.` at the beginning of URLs
		if config.App.SuppressWWW {
			providers.SuppressWWW(c)
		}
		// Force HTTPS protocol
		if config.App.ForceHTTPS {
			providers.ForceHTTPS(c)
		}
		// Move on the the next route
		c.Next()
	})

	// Use the Compression Middleware if enabled
	if config.Enabled["compression"] {
		app.Use(middleware.Compress(config.Compression))
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

	// Default Error Handler
	app.Settings.ErrorHandler = func(c *fiber.Ctx, err error) {
		code := fiber.StatusInternalServerError
		if e, ok := err.(*fiber.Error); ok {
			code = e.Code
		}
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		c.Status(code)
		c.JSON(fiber.Map{
			"success": false,
			"status":  code,
			"message": fmt.Sprintf("%s", err),
		})
	}

	// Start listening on the specified address
	err = app.Listen(config.App.Listen)
	if err != nil {
		// Exit the application
		exit(&config, app, err)
	}
}

func exit(config *configuration.Configuration, app *fiber.App, err error) {
	var dbErr error
	// Close database connection
	if config.Enabled["database"] {
		dbErr = err
		database.Close()
		if dbErr != nil {
			fmt.Printf("Closed database pool: %v\n", dbErr)
		} else {
			fmt.Println("Closed database pool")
		}
	}
	// Shutdown Fiber application
	var appErr error
	if err != nil {
		fmt.Printf("Shutdown Fiber application: %v", err)
		appErr = err
	} else {
		appErr = app.Shutdown()
		if appErr != nil {
			fmt.Printf("Shutdown Fiber application: %v", appErr)
		} else {
			fmt.Print("Shutdown Fiber application.")
		}
	}
	// Return with corresponding exit code
	if dbErr != nil || appErr != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
