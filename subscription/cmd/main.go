package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	subscriptionV1 "github.com/LearLocker/streaming/shared/pkg/openapi/subscription/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type SubscriptionStorage struct {
	mu            sync.RWMutex
	subscriptions map[string]*subscriptionV1.Subscription
}

func NewSubscriptionStorage() *SubscriptionStorage {
	return &SubscriptionStorage{
		subscriptions: make(map[string]*subscriptionV1.Subscription),
	}
}

func (s *SubscriptionStorage) GetSubscription(uuid string) *subscriptionV1.Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()

	subscription, ok := s.subscriptions[uuid]
	if !ok {
		return nil
	}

	return subscription
}

func (s *SubscriptionStorage) UpdateSubscription(uuid string, subscription *subscriptionV1.Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscriptions[uuid] = subscription
}

type SubscriptionHandler struct {
	storage *SubscriptionStorage
	catalog *clients.CatalogClient
}

func NewSubscriptionHandle(
	storage *SubscriptionStorage,
	catalog *clients.CatalogClient,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		storage: storage,
		catalog: catalog,
	}
}

func (h *SubscriptionHandler) CreateSubscription(
	ctx context.Context,
	req subscriptionV1.CreateSubscriptionRequest,
) (r subscriptionV1.CreateSubscriptionRes, _ error) {
	if req.PlanID == uuid.Nil {
		return &subscriptionV1.BadRequestError{
			Code:    404,
			Message: "plan_id is required",
		}, nil
	}

	plan, err := h.catalog.GetPlan(ctx, req.PlanID.String())
	if err != nil {
		if errors.Is(err, clients.ErrPlanNotFound) {
			return &api.NotFoundError{
				Message: "plan " + req.PlanID.String() + " not found",
			}, nil
		}
		// catalog недоступен — 500
		return &api.InternalServerError{
			Message: "failed to fetch plan: " + err.Error(),
		}, nil
	}

	newUUID := uuid.New()

	subscription := &subscriptionV1.Subscription{
		UUID:          subscriptionV1.NewOptUUID(newUUID),
		Status:        subscriptionV1.NewOptSubscriptionStatus(subscriptionV1.SubscriptionStatusPENDING),
		PlanID:        subscriptionV1.NewOptString(req.PlanID.String()),
		PaymentMethod: subscriptionV1.NewOptPaymentMethod(req.PaymentMethod),
		PlanName:      subscriptionV1.NewOptString(plan.Name),
		Amount:        subscriptionV1.NewOptInt64(plan.Price),
		Currency:      subscriptionV1.NewOptString(plan.Currency),
	}

	h.storage.UpdateSubscription(newUUID.String(), subscription)

	return subscription, nil
}

func (h *SubscriptionHandler) GetSubscriptionByUuid(
	ctx context.Context,
	params subscriptionV1.GetSubscriptionByUuidParams,
) (r subscriptionV1.GetSubscriptionByUuidRes, _ error) {
	subscription := h.storage.GetSubscription(params.UUID)
	if subscription == nil {
		return &subscriptionV1.NotFoundError{
			Code:    404,
			Message: "Subscription by uuid " + params.UUID + " not found",
		}, nil
	}

	return subscription, nil
}

func (h *SubscriptionHandler) PaySubscriptionByUuid(
	ctx context.Context,
	req subscriptionV1.PaySubscriptionRequest,
	params subscriptionV1.PaySubscriptionByUuidParams,
) (r subscriptionV1.PaySubscriptionByUuidRes, _ error) {
	subscription := h.storage.GetSubscription(params.UUID)
	if subscription == nil {
		return &subscriptionV1.NotFoundError{
			Code:    404,
			Message: "Subscription by uuid " + params.UUID + " not found",
		}, nil
	}

	if status, ok := subscription.Status.Get(); ok {
		switch status {
		case subscriptionV1.SubscriptionStatusACTIVE:
			return &subscriptionV1.ConflictError{Message: "subscription is already paid"}, nil
		case subscriptionV1.SubscriptionStatusCANCELLED:
			return &subscriptionV1.ConflictError{Message: "subscription is cancelled"}, nil
		}
	}

	// 3. получить тариф чтобы узнать duration_days
	planID, _ := sub.PlanID.Get()
	plan, err := h.catalog.GetPlan(ctx, planID)
	if err != nil {
		if errors.Is(err, clients.ErrPlanNotFound) {
			return &api.NotFoundError{
				Message: "plan " + planID + " not found",
			}, nil
		}
		return &api.InternalServerError{
			Message: "failed to fetch plan: " + err.Error(),
		}, nil
	}

	// 4. активировать — expires_at = now + duration_days из тарифа
	expiresAt := time.Now().UTC().AddDate(0, 0, int(plan.DurationDays))

	subscription.Status = subscriptionV1.NewOptSubscriptionStatus(subscriptionV1.SubscriptionStatusACTIVE)
	subscription.ExpiresAt = subscriptionV1.NewOptDateTime(expiresAt)

	h.storage.UpdateSubscription(params.UUID, subscription)

	return &subscriptionV1.PaySubscriptionResponse{
		UUID:          subscription.UUID,
		Status:        subscription.Status,
		PaymentMethod: subscription.PaymentMethod,
		ExpiresAt:     subscription.ExpiresAt,
	}, nil
}

func (h *SubscriptionHandler) CancelSubscriptionByUuid(
	ctx context.Context,
	params subscriptionV1.CancelSubscriptionByUuidParams,
) (r subscriptionV1.CancelSubscriptionByUuidRes, _ error) {
	subscription := h.storage.GetSubscription(params.UUID)
	if subscription == nil {
		return &subscriptionV1.NotFoundError{
			Code:    404,
			Message: "Subscription by uuid " + params.UUID + " not found",
		}, nil
	}

	if status, ok := subscription.Status.Get(); ok && status == subscriptionV1.SubscriptionStatusCANCELLED {
		return &subscriptionV1.ConflictError{
			Message: "subscription already cancelled",
		}, nil

	}

	subscription.Status = subscriptionV1.NewOptSubscriptionStatus(subscriptionV1.SubscriptionStatusCANCELLED)
	subscription.ExpiresAt = subscriptionV1.NewOptDateTime(time.Now())

	h.storage.UpdateSubscription(params.UUID, subscription)

	return &subscriptionV1.CancelSubscriptionByUuidNoContent{}, nil
}

func (h *SubscriptionHandler) NewError(ctx context.Context, err error) (r *subscriptionV1.GenericErrorStatusCode) {
	return &subscriptionV1.GenericErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response: subscriptionV1.GenericError{
			Code:    subscriptionV1.NewOptInt(http.StatusInternalServerError),
			Message: subscriptionV1.NewOptString(err.Error()),
		},
	}
}

func main() {
	catalogConn, err := grpc.NewClient(
		catalogAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to connect to catalog", slog.Any("err", err))
		os.Exit(1)
	}
	defer catalogConn.Close()

	catalogClient := clients.NewCatalogClient(catalogpb.NewCatalogServiceClient(catalogConn))
	// Создаем хранилище для данных
	storage := NewSubscriptionStorage()
	// Создаем обработчик API
	subscriptionHandler := NewSubscriptionHandle(storage, catalogClient)
	// Создаем OpenAPI сервер
	subscriptionServer, err := subscriptionV1.NewServer(subscriptionHandler)
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
