package payment

import (
	"context"
	"fmt"

	"github.com/LearLocker/streaming/payment/internal/model"
)

func (s *Service) ProcessPayment(ctx context.Context, info model.ProcessPaymentInfo) (*model.Payment, error) {
	if err := validateProcessInfo(info); err != nil {
		return nil, err
	}

	payment, err := s.paymentRepository.Create(ctx, info)
	if err != nil {
		return nil, fmt.Errorf("create payment: %w", err)
	}

	return payment, nil
}

func (s *Service) GetPayment(ctx context.Context, paymentID string) (*model.Payment, error) {
	if paymentID == "" {
		return nil, fmt.Errorf("payment_id is required")
	}

	payment, err := s.paymentRepository.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("get payment: %w", err)
	}

	return payment, nil
}

func (s *Service) RefundPayment(ctx context.Context, info model.RefundPaymentInfo) (*model.Payment, error) {
	if info.PaymentID == "" {
		return nil, fmt.Errorf("payment_id is required")
	}

	payment, err := s.paymentRepository.Refund(ctx, info)
	if err != nil {
		return nil, fmt.Errorf("refund payment: %w", err)
	}

	return payment, nil
}

func validateProcessInfo(info model.ProcessPaymentInfo) error {
	if info.SubscriptionID == "" {
		return fmt.Errorf("subscription_id is required")
	}
	if info.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if info.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if info.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if info.Method == "" {
		return fmt.Errorf("method is required")
	}
	return nil
}
