package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/Danila331/Swusher/internal/servers"
	"github.com/Danila331/Swusher/pkg/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Инициализация логгера
	// Используем zap для логирования
	// В данном случае используется development конфигурация для удобства отладки
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Fatal("failed to create logger", zap.Error(err))
	}
	defer logger.Sync()

	// Загрузка переменных окружения из файла .env
	err = godotenv.Load("./.env")
	if err != nil {
		logger.Fatal("failed to load env variables", zap.Error(err))
	}
	// Инициализация соединения с базой данных
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

	// Создание пула соединений с базой данных
	// Используем NexPgxPool для создания пула соединений
	pool, err := store.NexPgxPool(context.Background(), dbConfig)
	if err != nil {
		logger.Fatal("failed to create pgx pool", zap.Error(err))
	}
	defer pool.Close()

	// Создание таблиц, если они не существуют
	err = store.CreateTables(context.Background(), pool)
	if err != nil {
		logger.Fatal("failed to create tables", zap.Error(err))
	}

	// Настройка сервера метрик
	servers.SetupMetricsServer()

	// Запуск бекенд сервера
	servers.StartServer(logger, pool)
}
