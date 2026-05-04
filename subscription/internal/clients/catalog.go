package clients

import (
	"context"
	"errors"
	"fmt"

	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrPlanNotFound = errors.New("plan not found")

type Plan struct {
	PlanID       string
	Name         string
	Price        int64
	Currency     string
	DurationDays int32
	IsActive     bool
}

type CatalogClient struct {
	client catalogV1.CatalogServiceClient
}

func NewCatalogClient(client catalogV1.CatalogServiceClient) *CatalogClient {
	return &CatalogClient{client: client}
}

func (c *CatalogClient) GetPlan(ctx context.Context, planID string) (*Plan, error) {
	if planID == "" {
		return nil, fmt.Errorf("plan_id is required")
	}

	resp, err := c.client.GetPlan(ctx, &catalogV1.GetPlanRequest{
		PlanId: planID,
	})
	if err != nil {
		// переводим gRPC статус в доменную ошибку —
		// сервисный слой не должен знать про gRPC коды
		if status.Code(err) == codes.NotFound {
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("catalog.GetPlan: %w", err)
	}

	if resp.GetPlan() == nil {
		return nil, ErrPlanNotFound
	}

	return planFromProto(resp.GetPlan()), nil
}

func (c *CatalogClient) ListPlans(ctx context.Context, onlyActive bool) ([]*Plan, error) {
	resp, err := c.client.ListPlans(ctx, &catalogV1.ListPlansRequest{
		OnlyActive: onlyActive,
	})
	if err != nil {
		return nil, fmt.Errorf("catalog.ListPlans: %w", err)
	}

	plans := make([]*Plan, 0, len(resp.GetPlans()))
	for _, p := range resp.GetPlans() {
		plans = append(plans, planFromProto(p))
	}

	return plans, nil
}

// planFromProto — конвертер proto → локальная модель.
// Живёт здесь, а не в converter/, потому что это деталь клиента
func planFromProto(p *catalogV1.Plan) *Plan {
	return &Plan{
		PlanID:       p.GetPlanId(),
		Name:         p.GetName(),
		Price:        p.GetPrice(),
		Currency:     p.GetCurrency(),
		DurationDays: p.GetDurationDays(),
		IsActive:     p.GetIsActive(),
	}
}
