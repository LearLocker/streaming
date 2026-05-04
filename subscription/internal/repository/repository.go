package repository

import (
	"context"
	"time"

	"github.com/LearLocker/streaming/subscription/internal/model"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByUUID(ctx context.Context, uuid string) (*model.Subscription, error)
	UpdateStatus(ctx context.Context, uuid string, status model.SubscriptionStatus, expiresAt *time.Time) error
}
