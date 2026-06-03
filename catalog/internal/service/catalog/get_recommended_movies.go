package catalog

import (
	"context"
	"fmt"

	"github.com/LearLocker/streaming/catalog/internal/model"
)

func (s *Service) GetRecommendedMovies(ctx context.Context, userId string, limit int32) ([]*model.Movie, error) {
	if userId == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	movies, err := s.catalogRepository.GetRecommendedMovies(ctx, userId, limit)
	if err != nil {
		return nil, fmt.Errorf("get recommended movies: %w", err)
	}

	return movies, nil
}
