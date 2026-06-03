package repository

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/model"
)

type CatalogRepository interface {
	GetMovies(ctx context.Context, movieFilter model.GetMoviesFilter) ([]*model.Movie, int32, error)
	GetMovieById(ctx context.Context, imdbId string) (*model.Movie, error)
	AddMovie(ctx context.Context, info model.AddMovie) (string, error)
	UpdateReview(ctx context.Context, imdbId string, authorId string, text string, ranking model.Ranking) error
	GetRecommendedMovies(ctx context.Context, userId string, limit int32) ([]*model.Movie, error)
	GetGenres(ctx context.Context) ([]*model.Genre, error)
	GetPlan(ctx context.Context, planId string) (*model.Plan, error)
	ListPlans(ctx context.Context, onlyActive bool) ([]*model.Plan, error)
	GetUserFavouriteGenres(ctx context.Context, userID string) ([]string, error)
	GetRankings(ctx context.Context) ([]*model.Ranking, error)
}
