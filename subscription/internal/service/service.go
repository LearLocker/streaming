package service

import (
	"context"

	"github.com/LearLocker/streaming/subscription/internal/model"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, info model.CreateSubscriptionInfo) (*model.Subscription, error)
	GetSubscription(ctx context.Context, uuid string) (*model.Subscription, error)
	PaySubscription(ctx context.Context, info model.PaySubscriptionInfo) (*model.Subscription, error)
	CancelSubscription(ctx context.Context, uuid string) error
}
