package catalog

import (
	"context"
	"fmt"

	"github.com/LearLocker/streaming/catalog/internal/model"
)

func (s *Service) GetMovieById(ctx context.Context, imdbId string) (*model.Movie, error) {
	if imdbId == "" {
		return nil, fmt.Errorf("imdb_id is required")
	}

	movie, err := s.catalogRepository.GetMovieById(ctx, imdbId)
	if err != nil {
		return nil, fmt.Errorf("get movie: %w", err)
	}

	return movie, nil
}
