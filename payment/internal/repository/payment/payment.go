package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/LearLocker/streaming/payment/internal/model"
	"github.com/google/uuid"
)

func (r *repository) Create(ctx context.Context, info model.ProcessPaymentInfo) (*model.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p := &model.Payment{
		PaymentID:      uuid.New().String(),
		SubscriptionID: info.SubscriptionID,
		UserID:         info.UserID,
		Amount:         info.Amount,
		Currency:       info.Currency,
		Method:         info.Method,
		Status:         model.StatusSuccess,
		CreatedAt:      time.Now().UTC(),
	}

	r.payments[p.PaymentID] = p

	return p, nil
}

func (r *repository) GetByID(ctx context.Context, paymentID string) (*model.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.payments[paymentID]

	if !ok {
		return nil, fmt.Errorf("payment %s not found", paymentID)
	}

	return p, nil
}

func (r *repository) Refund(ctx context.Context, info model.RefundPaymentInfo) (*model.Payment, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p, ok := r.payments[info.PaymentID]
	if !ok {
		return nil, fmt.Errorf("payment %s not found", info.PaymentID)
	}

	if p.Status != model.StatusSuccess {
		return nil, fmt.Errorf("payment %s cannot be refunded: status is %s",
			info.PaymentID, p.Status)
	}

	now := time.Now().UTC()
	p.Status = model.StatusRefunded
	p.RefundedAt = &now
	p.Reason = info.Reason

	return p, nil
}
