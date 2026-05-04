package service

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/model"
)

type CatalogService interface {
	GetMovies(ctx context.Context, genreNames []string, page int32, pageSize int32) ([]*model.Movie, int32, error)
	GetMovieById(ctx context.Context, imdbId string) (*model.Movie, error)
	AddMovie(ctx context.Context, info model.AddMovie) (string, error)
	UpdateReview(ctx context.Context, imdbId string, authorId string, text string) (*model.UpdateReviewResult, error)
	GetRecommendedMovies(ctx context.Context, userId string, limit int32) ([]*model.Movie, error)
	GetGenres(ctx context.Context) ([]*model.Genre, error)
	GetPlan(ctx context.Context, planId string) (*model.Plan, error)
	ListPlans(ctx context.Context, onlyActive bool) ([]*model.Plan, error)
}
