package database

import (
	"context"
	"fmt"

	"github.com/mikeychowy/fiber-crayplate/app/configuration"

	"github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

// Instance of pgxpool
func Instance() *pgxpool.Pool {
	return pool
}

// Connect to the db through pool
func Connect(c context.Context, config *configuration.DatabaseConfiguration) (err error) {
	pool, err = pgxpool.Connect(*&c, "host="+config.Host+" port="+fmt.Sprintf("%d", config.Port)+" user="+config.Username+" dbname="+config.Database+" password="+config.Password)
	if err != nil {
		fmt.Printf("Error Connecting to Database! Reason: %s", err)
	}
	return err
}

// Close pool connection
func Close() {
	pool.Close()
}
