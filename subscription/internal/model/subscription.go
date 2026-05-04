package model

import "time"

type SubscriptionStatus string

const (
	StatusPending         SubscriptionStatus = "PENDING"
	StatusAwaitingPayment SubscriptionStatus = "AWAITING_PAYMENT"
	StatusActive          SubscriptionStatus = "ACTIVE"
	StatusExpired         SubscriptionStatus = "EXPIRED"
	StatusCancelled       SubscriptionStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodCard   PaymentMethod = "CARD"
	PaymentMethodSBP    PaymentMethod = "SBP"
	PaymentMethodWallet PaymentMethod = "WALLET"
)

type Subscription struct {
	UUID          string
	PlanID        string
	PlanName      string
	Amount        int64
	Currency      string
	PaymentMethod PaymentMethod
	Status        SubscriptionStatus
	ExpiresAt     *time.Time
}

// CreateSubscriptionInfo — данные от клиента при создании
type CreateSubscriptionInfo struct {
	PlanID        string
	PaymentMethod PaymentMethod
}

// PaySubscriptionInfo — данные при оплате
type PaySubscriptionInfo struct {
	UUID          string
	PaymentMethod PaymentMethod
}
