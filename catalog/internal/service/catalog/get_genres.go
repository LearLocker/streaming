package catalog

import (
	"context"
	"fmt"

	"github.com/LearLocker/streaming/catalog/internal/model"
)

func (s *Service) GetGenres(ctx context.Context) ([]*model.Genre, error) {
	genres, err := s.catalogRepository.GetGenres(ctx)
	if err != nil {
		return nil, fmt.Errorf("get genres: %w", err)
	}

	return genres, nil
}
