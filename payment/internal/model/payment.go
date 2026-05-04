package model

import "time"

type PaymentStatus string

const (
	StatusPending    PaymentStatus = "PAYMENT_STATUS_PENDING"
	StatusProcessing PaymentStatus = "PAYMENT_STATUS_PROCESSING"
	StatusSuccess    PaymentStatus = "PAYMENT_STATUS_SUCCESS"
	StatusFailed     PaymentStatus = "PAYMENT_STATUS_FAILED"
	StatusRefunded   PaymentStatus = "PAYMENT_STATUS_REFUNDED"
)

// Payment — доменная модель, не зависит от proto или БД
type Payment struct {
	PaymentID      string
	SubscriptionID string
	UserID         string
	Amount         int64
	Currency       string
	Method         string
	Status         PaymentStatus
	Reason         string
	CreatedAt      time.Time
	RefundedAt     *time.Time
}

// ProcessPaymentInfo — данные для создания платежа
type ProcessPaymentInfo struct {
	SubscriptionID string
	UserID         string
	Amount         int64
	Currency       string
	Method         string
}

// RefundPaymentInfo — данные для возврата
type RefundPaymentInfo struct {
	PaymentID string
	Reason    string
}
