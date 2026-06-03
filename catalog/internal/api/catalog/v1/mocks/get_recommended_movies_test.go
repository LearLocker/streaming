package v1

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *APISuite) TestGetRecommendedMoviesSuccess() {
	var (
		userID = gofakeit.UUID()
		limit  = int32(gofakeit.IntRange(5, 20))

		req = &catalogV1.GetRecommendedMoviesRequest{
			UserId: userID,
			Limit:  limit,
		}

		expectedMovies = []*model.Movie{
			{
				ImdbID:     gofakeit.UUID(),
				Title:      gofakeit.MovieName(),
				PosterPath: gofakeit.ImageURL(1920, 1080),
				Ranking: model.Ranking{
					RankingValue: gofakeit.IntRange(1, 100),
					RankingName:  gofakeit.Word(),
				},
			},
			{
				ImdbID:     gofakeit.UUID(),
				Title:      gofakeit.MovieName(),
				PosterPath: gofakeit.ImageURL(1920, 1080),
			},
		}
	)

	s.catalogService.On("GetRecommendedMovies", s.ctx, userID, limit).Return(expectedMovies, nil)

	res, err := s.api.GetRecommendedMovies(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.GetMovies(), len(expectedMovies))
}

func (s *APISuite) TestGetRecommendedMoviesEmptyResult() {
	var (
		userID = gofakeit.UUID()
		limit  = int32(10)

		req = &catalogV1.GetRecommendedMoviesRequest{
			UserId: userID,
			Limit:  limit,
		}
	)

	s.catalogService.On("GetRecommendedMovies", s.ctx, userID, limit).Return([]*model.Movie{}, nil)

	res, err := s.api.GetRecommendedMovies(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Empty(res.GetMovies())
}

func (s *APISuite) TestGetRecommendedMoviesServiceError() {
	var (
		serviceErr = gofakeit.Error()
		userID     = gofakeit.UUID()
		limit      = int32(10)

		req = &catalogV1.GetRecommendedMoviesRequest{
			UserId: userID,
			Limit:  limit,
		}
	)

	s.catalogService.On("GetRecommendedMovies", s.ctx, userID, limit).Return(nil, serviceErr)

	res, err := s.api.GetRecommendedMovies(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
