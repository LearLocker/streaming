package catalog

import (
	"context"
	"fmt"

	"github.com/LearLocker/streaming/catalog/internal/model"
)

func (s *Service) GetPlan(ctx context.Context, planId string) (*model.Plan, error) {
	if planId == "" {
		return nil, fmt.Errorf("plan_id is required")
	}

	plan, err := s.catalogRepository.GetPlan(ctx, planId)
	if err != nil {
		return nil, fmt.Errorf("get plan: %w", err)
	}

	return plan, nil
}
