package catalog

import (
	"context"
	"fmt"

	"github.com/LearLocker/streaming/catalog/internal/model"
)

func (s *Service) ListPlans(ctx context.Context, onlyActive bool) ([]*model.Plan, error) {
	plans, err := s.catalogRepository.ListPlans(ctx, onlyActive)
	if err != nil {
		return nil, fmt.Errorf("list plans: %w", err)
	}

	return plans, nil
}
