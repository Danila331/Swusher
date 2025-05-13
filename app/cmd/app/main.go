package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Danila331/ShareHub/internal/servers"
	"github.com/Danila331/ShareHub/pkg/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Initialize the logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("failed to create logger", zap.Error(err))
	}
	defer logger.Sync() // flushes buffer, if any

	// load env variables
	err = godotenv.Load("./.env")
	if err != nil {
		logger.Fatal("failed to load env variables", zap.Error(err))
	}
	// Initialize the database connection
	port, _ := strconv.Atoi(os.Getenv("POSTGRESQL_PORT"))
	dbConfig := store.NewDBConfig(os.Getenv("POSTGRESQL_HOST"),
		port,
		os.Getenv("POSTGRESQL_USER"),
		os.Getenv("POSTGRESQL_PASSWORD"),
		os.Getenv("POSTGRESQL_DBNAME"),
		"disable",
		os.Getenv("POSTGRESQL_CONNECTION_STRING"),
		10,
		2,
		5*time.Minute,
		30*time.Minute,
	)

	// Create a new pgx pool
	pool, err := store.NexPgxPool(context.Background(), dbConfig)
	if err != nil {
		logger.Fatal("failed to create pgx pool", zap.Error(err))
	}
	defer pool.Close()

	// Create Tables if not exist
	err = store.CreateTables(context.Background(), pool)
	if err != nil {
		logger.Fatal("failed to create tables", zap.Error(err))
	}

	// Start the server
	servers.StartServer(logger)
}
