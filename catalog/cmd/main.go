package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	catalogApiV1 "github.com/LearLocker/streaming/catalog/internal/api/catalog/v1"
	catalogRepository "github.com/LearLocker/streaming/catalog/internal/repository/catalog"
	catalogService "github.com/LearLocker/streaming/catalog/internal/service/catalog"
	"github.com/LearLocker/streaming/catalog/internal/storage/mongoDb"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
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

	mongoClient := mongoDb.Connect()

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	repo := catalogRepository.NewRepository(mongoClient)
	service := catalogService.NewService(repo)
	api := catalogApiV1.NewAPI(service)

	// Регистрируем наш сервис
	catalogV1.RegisterCatalogServiceServer(grpcServer, api)

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = grpcServer.Serve(lis)
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
	grpcServer.GracefulStop()
	log.Println("✅ Server stopped")
}
