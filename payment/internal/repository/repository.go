package repository

import (
	"context"

	"github.com/LearLocker/streaming/payment/internal/model"
)

type PaymentRepository interface {
	Create(ctx context.Context, info model.ProcessPaymentInfo) (*model.Payment, error)
	GetByID(ctx context.Context, paymentID string) (*model.Payment, error)
	Refund(ctx context.Context, info model.RefundPaymentInfo) (*model.Payment, error)
}
