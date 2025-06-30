package main

import (
	"context"
	"log"
	"net"
	pb "notification_service/api/protobuf/notification"
	"notification_service/internal/database"
	"notification_service/internal/handlers"
	"notification_service/internal/queue"
	"notification_service/internal/service"
	"notification_service/internal/workers"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedNotificationServiceServer
	handler handlers.NotificationHandler
}

func (s *Server) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	notificationID, err := s.handler.SendNotification(ctx, req.UserId, req.Type, req.Message, req.Priority, req.ScheduledAt)
	if err != nil {
		log.Println("Error sending notification: ", err)
		return nil, err
	}

	return &pb.SendNotificationResponse{NotificationId: notificationID}, nil
}

func (s *Server) GetStatus(ctx context.Context, req *pb.GetStatusRequest) (*pb.GetStatusResponse, error) {
	notification, err := s.handler.GetStatus(ctx, req.NotificationId)
	if err != nil {
		log.Println("Error getting status: ", err)
		return nil, err
	}

	var deliveredAt int64
	if !notification.DeliveredAt.IsZero() {
		deliveredAt = notification.DeliveredAt.Unix()
	}

	return &pb.GetStatusResponse{
		Status:      string(notification.Status),
		Attempts:    notification.Attempts,
		LastError:   notification.LastError,
		DeliveredAt: deliveredAt,
	}, nil
}

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return
	}
	queue, err := queue.NewNatsQueue(nats.DefaultURL)
	if err != nil {
		log.Printf("Error connecting to queue: %v", err)
		return
	}
	notificationService := service.NewNotificationsService(db)
	notificationHandler := handlers.NewNotificationHandler(notificationService, queue)

	workerPool := workers.NewPool(5, queue, notificationService)
	workerPool.Start()
	defer workerPool.Stop()

	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Printf("Failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, &Server{handler: *notificationHandler})

	log.Println("Server started on port 50051")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
