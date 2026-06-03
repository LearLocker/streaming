package catalog

import (
	"context"
	"fmt"

	"github.com/LearLocker/streaming/catalog/internal/model"
)

func (s *Service) GetMovies(ctx context.Context, movieFilter model.GetMoviesFilter) ([]*model.Movie, int32, error) {
	movies, total, err := s.catalogRepository.GetMovies(ctx, movieFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("get movies: %w", err)
	}

	return movies, total, nil
}
