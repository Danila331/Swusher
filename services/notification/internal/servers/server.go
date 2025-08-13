package servers

import (
	"log"
	"net"

	"github.com/Danila331/Swusher/notification-server/internal/handlers"
	notification "github.com/Danila331/Swusher/notification-server/internal/pb/notification/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func StartServer(logger *zap.Logger) {
	logger.Info("Starting Notification Service...")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	notification.RegisterNotificationServiceServer(grpcServer, handlers.NewNotificationServiceServer())

	log.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
