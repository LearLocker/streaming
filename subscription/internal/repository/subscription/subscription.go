package subscription

import (
	"context"
	"fmt"
	"time"

	"github.com/LearLocker/streaming/subscription/internal/model"
)

func (r *repository) Create(
	_ context.Context,
	sub *model.Subscription,
) error {

	r.mu.Lock()
	defer r.mu.Unlock()
	r.subscriptions[sub.UUID] = sub

	return nil
}

func (r *repository) GetByUUID(
	_ context.Context,
	uuid string,
) (*model.Subscription, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	sub, ok := r.subscriptions[uuid]
	if !ok {
		return nil, fmt.Errorf("subscription %s not found", uuid)
	}

	return sub, nil
}

func (r *repository) UpdateStatus(
	_ context.Context,
	uuid string,
	status model.SubscriptionStatus,
	expiresAt *time.Time,
) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	sub, ok := r.subscriptions[uuid]
	if !ok {
		return fmt.Errorf("subscription %s not found", uuid)
	}

	sub.Status = status
	sub.ExpiresAt = expiresAt

	return nil
}
