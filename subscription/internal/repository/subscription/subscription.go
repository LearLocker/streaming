package subscription

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/LearLocker/streaming/subscription/internal/model"
)

func (r *repository) Create(
	ctx context.Context,
	sub *model.Subscription,
) error {

	const query = `
        INSERT INTO subscriptions (
            uuid, plan_id, plan_name, amount,
            currency, payment_method, status
        ) VALUES (
            $1, $2, $3, $4,
            $5, $6, $7
        )
    `
	_, err := r.db.ExecContext(ctx, query,
		sub.UUID,
		sub.PlanID,
		sub.PlanName,
		sub.Amount,
		sub.Currency,
		string(sub.PaymentMethod),
		string(sub.Status),
	)
	if err != nil {
		return fmt.Errorf("insert subscription: %w", err)
	}

	return nil
}

func (r *repository) GetByUUID(
	ctx context.Context,
	uuid string,
) (*model.Subscription, error) {

	const query = `
        SELECT uuid, plan_id, plan_name, amount,
               currency, payment_method, status, expires_at
        FROM subscriptions
        WHERE uuid = $1
    `

	var sub model.Subscription
	var expiresAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&sub.UUID,
		&sub.PlanID,
		&sub.PlanName,
		&sub.Amount,
		&sub.Currency,
		&sub.PaymentMethod,
		&sub.Status,
		&expiresAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("subscription %s not found", uuid)
	}
	if err != nil {
		return nil, fmt.Errorf("get subscription: %w", err)
	}

	if expiresAt.Valid {
		sub.ExpiresAt = &expiresAt.Time
	}

	return &sub, nil
}

func (r *repository) UpdateStatus(
	ctx context.Context,
	uuid string,
	status model.SubscriptionStatus,
	expiresAt *time.Time,
) error {

	const query = `
        UPDATE subscriptions
        SET status     = $1,
            expires_at = $2,
            updated_at = now()
        WHERE uuid = $3
    `

	res, err := r.db.ExecContext(ctx, query, string(status), expiresAt, uuid)
	if err != nil {
		return fmt.Errorf("update status: %w", err)
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return fmt.Errorf("subscription %s not found", uuid)
	}

	return nil
}
