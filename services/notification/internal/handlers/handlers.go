package handlers

import (
	"context"
	"fmt"

	notification "github.com/Danila331/Swusher/notification-server/pkg/pb/notification/proto"
)

type NotificationServiceServerImpl struct {
	notification.UnimplementedNotificationServiceServer
}

func NewNotificationServiceServer() notification.NotificationServiceServer {
	return &NotificationServiceServerImpl{}
}

// SendNotification отправляет уведомление пользователю
func (s *NotificationServiceServerImpl) SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	// Пример: просто логируем и возвращаем успех
	fmt.Printf("SendNotification: user_id=%s, message=%s\n", req.GetUserId(), req.GetMessage())
	return &notification.SendNotificationResponse{
		Success: true,
		Error:   "",
	}, nil
}

func (s *NotificationServiceServerImpl) GetUniqueCode(ctx context.Context, req *notification.GetUniqueCodeRequest) (*notification.GetUniqueCodeResponse, error) {
	// Пример: генерируем уникальный код (можно заменить на свою логику)
	code := "CODE-" + req.GetUserId()
	return &notification.GetUniqueCodeResponse{
		UniqueCode: code,
		Success:    true,
		Error:      "",
	}, nil
}
