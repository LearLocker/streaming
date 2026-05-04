package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	paymentApiV1 "github.com/LearLocker/streaming/payment/internal/api/payment/v1"
	paymentRepository "github.com/LearLocker/streaming/payment/internal/repository/payment"
	paymentService "github.com/LearLocker/streaming/payment/internal/service/payment"
	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50051

func main() {

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	repo := paymentRepository.NewRepository()
	service := paymentService.NewService(repo)
	api := paymentApiV1.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
