package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LearLocker/streaming/subscription/internal/migrator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	subscriptionV1 "github.com/LearLocker/streaming/shared/pkg/openapi/subscription/v1"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	subscriptionApiV1 "github.com/LearLocker/streaming/subscription/internal/api/subscription/v1"
	"github.com/LearLocker/streaming/subscription/internal/clients"
	subscriptionRepository "github.com/LearLocker/streaming/subscription/internal/repository/subscription"
	subscriptionService "github.com/LearLocker/streaming/subscription/internal/service/subscription"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("failed to load .env file: %v\n", err)
	}

	dbURI := os.Getenv("DB_URI")
	if dbURI == "" {
		log.Fatal("DB_URI is required")
	}

	// открываем пул соединений сразу через stdlib
	// pgx используется как драйвер для database/sql
	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		log.Fatalf("failed to open db: %v\n", err)
	}
	defer db.Close()

	// настройка пула
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// проверяем соединение
	if err = db.PingContext(ctx); err != nil {
		log.Fatalf("database unavailable: %v\n", err)
	}

	log.Println("database connected")

	// миграции
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "migrations"
	}

	migratorRunner := migrator.NewMigrator(db, migrationsDir)
	if err = migratorRunner.Up(); err != nil {
		log.Fatalf("migration failed: %v\n", err)
	}

	log.Println("migrations applied")

	catalogAddr := os.Getenv("CATALOG_ADDR")
	if catalogAddr == "" {
		catalogAddr = "localhost:50052"
	}

	catalogConn, err := grpc.NewClient(
		catalogAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect to catalog: %v\n", err)
		os.Exit(1)
	}
	defer catalogConn.Close()

	catalogClient := clients.NewCatalogClient(catalogV1.NewCatalogServiceClient(catalogConn))
	subRepo := subscriptionRepository.NewRepository(db)
	subService := subscriptionService.NewService(subRepo, catalogClient)
	subAPI := subscriptionApiV1.NewAPI(subService)

	// Создаем OpenAPI сервер
	subscriptionServer, err := subscriptionV1.NewServer(subAPI)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	// Добавляем middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	// Монтируем обработчики OpenAPI
	r.Mount("/", subscriptionServer)

	// Запускаем HTTP-сервер
	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout, // Защита от Slowloris атак - тип DDoS-атаки, при которой
		// атакующий умышленно медленно отправляет HTTP-заголовки, удерживая соединения открытыми и истощая
		// пул доступных соединений на сервере. ReadHeaderTimeout принудительно закрывает соединение,
		// если клиент не успел отправить все заголовки за отведенное время.
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
