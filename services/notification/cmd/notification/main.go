package main

import (
	"github.com/Danila331/Swusher/notification-server/internal/servers"
	"go.uber.org/zap"
)

func main() {
	// Инициализация gRPC сервера
	logger := zap.NewExample()
	defer logger.Sync() // Запланировать синхронизацию логов перед выходом
	logger.Info("Starting notification service...")

	servers.StartServer(logger) // Запуск gRPC сервера с логгером
}
