package v1

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *APISuite) TestGetMoviesSuccess() {
	var (
		genreNames = []string{
			gofakeit.MovieGenre(),
			gofakeit.MovieGenre(),
		}
		page     = int32(gofakeit.IntRange(1, 10))
		pageSize = int32(gofakeit.IntRange(5, 50))

		req = &catalogV1.GetMoviesRequest{
			GenreNames: genreNames,
			Page:       page,
			PageSize:   pageSize,
		}

		// Повторяем ту же логику, что и в хендлере,
		// чтобы мок сработал на идентичный аргумент
		expectedFilter = model.GetMoviesFilter{
			GenreNames: genreNames,
			Page:       page,
			PageSize:   pageSize,
		}

		expectedMovies = []*model.Movie{
			{
				ImdbID:     gofakeit.UUID(),
				Title:      gofakeit.MovieName(),
				PosterPath: gofakeit.ImageURL(1920, 1080),
				YouTubeID:  "dQw4w9WgXcQ",
			},
			{
				ImdbID:     gofakeit.UUID(),
				Title:      gofakeit.MovieName(),
				PosterPath: gofakeit.ImageURL(1920, 1080),
				YouTubeID:  "dQw4w9WgXcQ",
			},
		}
	)

	s.catalogService.On("GetMovies", s.ctx, expectedFilter).Return(expectedMovies, 2, nil)

	res, err := s.api.GetMovies(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.GetMovies(), len(expectedMovies))
}

func (s *APISuite) TestGetMoviesNoGenreFilter() {
	var (
		page     = int32(1)
		pageSize = int32(10)

		req = &catalogV1.GetMoviesRequest{
			Page:     page,
			PageSize: pageSize,
		}

		expectedFilter = model.GetMoviesFilter{
			GenreNames: nil,
			Page:       page,
			PageSize:   pageSize,
		}

		expectedMovies = []*model.Movie{
			{ImdbID: gofakeit.UUID(), Title: gofakeit.MovieName()},
			{ImdbID: gofakeit.UUID(), Title: gofakeit.MovieName()},
		}
	)

	s.catalogService.On("GetMovies", s.ctx, expectedFilter).Return(expectedMovies, 2, nil)

	res, err := s.api.GetMovies(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.GetMovies(), len(expectedMovies))
}

func (s *APISuite) TestGetMoviesEmptyResult() {
	var (
		req = &catalogV1.GetMoviesRequest{
			Page:     1,
			PageSize: 10,
		}

		expectedFilter = model.GetMoviesFilter{
			Page:     1,
			PageSize: 10,
		}
	)

	s.catalogService.On("GetMovies", s.ctx, expectedFilter).Return([]*model.Movie{}, 0, nil)

	res, err := s.api.GetMovies(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Empty(res.GetMovies())
}

func (s *APISuite) TestGetMoviesServiceError() {
	var (
		serviceErr = gofakeit.Error()

		req = &catalogV1.GetMoviesRequest{
			Page:     1,
			PageSize: 10,
		}

		expectedFilter = model.GetMoviesFilter{
			Page:     1,
			PageSize: 10,
		}
	)

	s.catalogService.On("GetMovies", s.ctx, expectedFilter).Return(nil, 0, serviceErr)

	res, err := s.api.GetMovies(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
