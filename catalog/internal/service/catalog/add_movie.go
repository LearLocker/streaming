package catalog

import (
	"context"
	"fmt"

	"github.com/LearLocker/streaming/catalog/internal/model"
)

func (s *Service) AddMovie(ctx context.Context, info model.AddMovie) (string, error) {
	if info.ImdbID == "" || info.Title == "" {
		return "", fmt.Errorf("imdb_id and title are required")
	}

	insertedID, err := s.catalogRepository.AddMovie(ctx, info)
	if err != nil {
		return "", fmt.Errorf("add movie: %w", err)
	}

	return insertedID, nil
}
