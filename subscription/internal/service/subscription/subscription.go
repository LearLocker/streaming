package subscription

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LearLocker/streaming/subscription/internal/clients"
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/google/uuid"
)

func (s *service) CreateSubscription(
	ctx context.Context,
	info model.CreateSubscriptionInfo,
) (*model.Subscription, error) {

	if info.PlanID == "" {
		return nil, fmt.Errorf("plan_id is required")
	}

	// получить данные тарифа из Catalog Service
	plan, err := s.catalogClient.GetPlan(ctx, info.PlanID)
	if err != nil {
		if errors.Is(err, clients.ErrPlanNotFound) {
			return nil, fmt.Errorf("plan %s not found", info.PlanID)
		}
		return nil, fmt.Errorf("get plan: %w", err)
	}

	sub := &model.Subscription{
		UUID:          uuid.New().String(),
		PlanID:        plan.PlanID,
		PlanName:      plan.Name,
		Amount:        plan.Price,
		Currency:      plan.Currency,
		PaymentMethod: info.PaymentMethod,
		Status:        model.StatusPending,
	}

	if err = s.subscriptionRepository.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("create subscription: %w", err)
	}

	return sub, nil
}

func (s *service) GetSubscription(
	ctx context.Context,
	uuid string,
) (*model.Subscription, error) {

	if uuid == "" {
		return nil, fmt.Errorf("uuid is required")
	}

	sub, err := s.subscriptionRepository.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	return sub, nil
}

func (s *service) PaySubscription(
	ctx context.Context,
	info model.PaySubscriptionInfo,
) (*model.Subscription, error) {

	sub, err := s.subscriptionRepository.GetByUUID(ctx, info.UUID)
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	switch sub.Status {
	case model.StatusActive:
		return nil, fmt.Errorf("subscription is already paid")
	case model.StatusCancelled:
		return nil, fmt.Errorf("subscription is cancelled")
	}

	// получить duration_days из тарифа для расчёта expires_at
	plan, err := s.catalogClient.GetPlan(ctx, sub.PlanID)
	if err != nil {
		return nil, fmt.Errorf("get plan: %w", err)
	}

	expiresAt := time.Now().UTC().AddDate(0, 0, int(plan.DurationDays))

	if err = s.subscriptionRepository.UpdateStatus(
		ctx, sub.UUID, model.StatusActive, &expiresAt,
	); err != nil {
		return nil, fmt.Errorf("activate subscription: %w", err)
	}

	sub.Status = model.StatusActive
	sub.ExpiresAt = &expiresAt

	return sub, nil
}

func (s *service) CancelSubscription(
	ctx context.Context,
	uuid string,
) error {

	sub, err := s.subscriptionRepository.GetByUUID(ctx, uuid)
	if err != nil {
		return fmt.Errorf("get subscription: %w", err)
	}

	if sub.Status == model.StatusCancelled {
		return fmt.Errorf("subscription is already cancelled")
	}

	if err = s.subscriptionRepository.UpdateStatus(
		ctx, uuid, model.StatusCancelled, nil,
	); err != nil {
		return fmt.Errorf("cancel subscription: %w", err)
	}

	return nil
}
