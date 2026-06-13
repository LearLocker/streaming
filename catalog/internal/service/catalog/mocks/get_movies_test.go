package catalog

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestGetMoviesSuccess() {
	var (
		filter = model.GetMoviesFilter{
			GenreNames: []string{gofakeit.MovieGenre(), gofakeit.MovieGenre()},
			Page:       int32(gofakeit.IntRange(1, 10)),
			PageSize:   int32(gofakeit.IntRange(5, 50)),
		}

		expectedMovies = []*model.Movie{
			{
				ImdbID:     gofakeit.UUID(),
				Title:      gofakeit.MovieName(),
				PosterPath: gofakeit.ImageURL(1920, 1080),
				YouTubeID:  "dQw4w9WgXcQ",
				Genre: []model.Genre{{
					GenreID:   gofakeit.IntRange(1, 20),
					GenreName: filter.GenreNames[0],
				}},
			},
		}
	)

	s.catalogRepository.On("GetMovies", s.ctx, filter).Return(expectedMovies, 1, nil)

	movies, total, err := s.service.GetMovies(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().NotNil(movies)
	s.Require().Len(movies, len(expectedMovies))
	s.Require().Equal(total, 1)
	s.Require().Equal(expectedMovies, movies)
}

func (s *ServiceSuite) TestGetMoviesNoGenreFilter() {
	var (
		filter = model.GetMoviesFilter{
			Page:     1,
			PageSize: 10,
		}

		expectedMovies = []*model.Movie{
			{ImdbID: gofakeit.UUID(), Title: gofakeit.MovieName()},
			{ImdbID: gofakeit.UUID(), Title: gofakeit.MovieName()},
		}
	)

	s.catalogRepository.On("GetMovies", s.ctx, filter).Return(expectedMovies, nil)

	movies, total, err := s.service.GetMovies(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Len(movies, 2)
	s.Require().Equal(total, 2)
}

func (s *ServiceSuite) TestGetMoviesEmptyResult() {
	var (
		filter = model.GetMoviesFilter{
			Page:     1,
			PageSize: 10,
		}
	)

	s.catalogRepository.On("GetMovies", s.ctx, filter).Return([]*model.Movie{}, nil)

	movies, total, err := s.service.GetMovies(s.ctx, filter)
	s.Require().NoError(err)
	s.Require().Empty(movies)
	s.Require().Equal(total, 0)
}

func (s *ServiceSuite) TestGetMoviesRepoError() {
	var (
		repoErr = gofakeit.Error()
		filter  = model.GetMoviesFilter{
			Page:     1,
			PageSize: 10,
		}
	)

	s.catalogRepository.On("GetMovies", s.ctx, filter).Return(nil, repoErr)

	movies, total, err := s.service.GetMovies(s.ctx, filter)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(movies)
	s.Require().Equal(total, 0)
}
