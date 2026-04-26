package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type paymentService struct {
	paymentV1.UnimplementedPaymentServiceServer

	payments map[string]*paymentV1.Payment
	mu       sync.RWMutex
}

func (s *paymentService) ProcessPayment(
	_ context.Context, req *paymentV1.ProcessPaymentRequest,
) (*paymentV1.ProcessPaymentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := &paymentV1.Payment{
		PaymentId:      uuid.New().String(),
		SubscriptionId: req.GetSubscriptionId(),
		UserId:         req.GetUserId(),
		Amount:         req.GetAmount(),
		Currency:       req.GetCurrency(),
		Method:         req.GetMethod(),
		Status:         paymentV1.PaymentStatus_PAYMENT_STATUS_SUCCESS,
		CreatedAt:      timestamppb.New(time.Now()),
	}

	s.payments[p.PaymentId] = p

	return &paymentV1.ProcessPaymentResponse{
		PaymentId: p.PaymentId,
		Status:    p.Status,
		Message:   "payment processed successfully",
	}, nil
}
func (s *paymentService) GetPayment(
	_ context.Context,
	req *paymentV1.GetPaymentRequest,
) (*paymentV1.GetPaymentResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if req.GetPaymentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	payment, ok := s.payments[req.GetPaymentId()]

	if !ok {
		return nil, status.Errorf(codes.NotFound,
			"payment %s not found", req.GetPaymentId())
	}

	return &paymentV1.GetPaymentResponse{
		Payment: payment,
	}, nil
}

func (s *paymentService) RefundPayment(
	_ context.Context,
	req *paymentV1.RefundPaymentRequest,
) (*paymentV1.RefundPaymentResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if req.GetPaymentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	p, ok := s.payments[req.GetPaymentId()]
	if !ok {
		return nil, status.Errorf(codes.NotFound,
			"payment %s not found", req.GetPaymentId())
	}

	switch p.Status {
	case paymentV1.PaymentStatus_PAYMENT_STATUS_SUCCESS:
		// единственный допустимый статус для возврата

	case paymentV1.PaymentStatus_PAYMENT_STATUS_REFUNDED:
		return nil, status.Errorf(codes.FailedPrecondition,
			"payment %s is already refunded", req.GetPaymentId())

	case paymentV1.PaymentStatus_PAYMENT_STATUS_FAILED:
		return nil, status.Errorf(codes.FailedPrecondition,
			"payment %s failed, cannot refund", req.GetPaymentId())

	default:
		return nil, status.Errorf(codes.FailedPrecondition,
			"payment %s has status %s, cannot refund",
			req.GetPaymentId(), p.Status.String())
	}

	now := time.Now()
	p.Status = paymentV1.PaymentStatus_PAYMENT_STATUS_REFUNDED
	p.RefundedAt = timestamppb.New(now)
	p.Reason = req.GetReason()

	return &paymentV1.RefundPaymentResponse{
		PaymentId:  p.PaymentId,
		Status:     p.Status,
		RefundedAt: timestamppb.New(now),
	}, nil
}

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
	service := &paymentService{
		payments: make(map[string]*paymentV1.Payment),
	}

	paymentV1.RegisterPaymentServiceServer(s, service)

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
