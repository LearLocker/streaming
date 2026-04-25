package main

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
)

var ErrPlanNotFound = errors.New("plan not found")

type Plan struct {
	PlanID       string
	Name         string
	Price        int64
	Currency     string
	DurationDays int32
}

type CatalogClient struct {
	client catalogV1.CatalogServiceClient
}

func NewCatalogClient(client catalogV1.CatalogServiceClient) *CatalogClient {
	return &CatalogClient{client: client}
}

func (c *CatalogClient) GetPlan(ctx context.Context, planID string) (*Plan, error) {
	resp, err := c.client.GetPlan(ctx, &catalogV1.GetPlanRequest{
		PlanId: planID,
	})
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrPlanNotFound
		}
		return nil, fmt.Errorf("catalog.GetPlan: %w", err)
	}

	return &Plan{
		PlanID:       resp.Plan.PlanId,
		Name:         resp.Plan.Name,
		Price:        resp.Plan.Price,
		Currency:     resp.Plan.Currency,
		DurationDays: resp.Plan.DurationDays,
	}, nil
}
