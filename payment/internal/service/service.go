package service

import (
	"context"

	"github.com/LearLocker/streaming/payment/internal/model"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, info model.ProcessPaymentInfo) (*model.Payment, error)
	GetPayment(ctx context.Context, paymentID string) (*model.Payment, error)
	RefundPayment(ctx context.Context, info model.RefundPaymentInfo) (*model.Payment, error)
}
