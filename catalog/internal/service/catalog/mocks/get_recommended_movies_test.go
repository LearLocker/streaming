package catalog

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestGetRecommendedMoviesSuccess() {
	var (
		UserID = gofakeit.UUID()
		Limit  = int32(gofakeit.IntRange(5, 20))

		// GetRecommendedMovies внутри запрашивает жанры пользователя,
		// поэтому мокаем репозиторий на оба вызова
		userGenres = []model.Genre{
			{GenreID: gofakeit.IntRange(1, 20), GenreName: gofakeit.MovieGenre()},
			{GenreID: gofakeit.IntRange(1, 20), GenreName: gofakeit.MovieGenre()},
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

	s.catalogRepository.On("GetUsersFavouriteGenres", s.ctx, UserID).Return(userGenres, nil)
	s.catalogRepository.On("GetRecommendedMovies", s.ctx, userGenres, Limit).Return(expectedMovies, nil)

	movies, err := s.catalogRepository.GetRecommendedMovies(s.ctx, UserID, Limit)
	s.Require().NoError(err)
	s.Require().NotNil(movies)
	s.Require().Len(movies, len(expectedMovies))
	s.Require().Equal(expectedMovies, movies)
}

func (s *ServiceSuite) TestGetRecommendedMoviesGetGenresError() {
	var (
		repoErr = gofakeit.Error()
		UserID  = gofakeit.UUID()
		Limit   = int32(10)
	)

	// Ошибка на первом шаге — за жанрами идти не должен
	s.catalogRepository.On("GetUsersFavouriteGenres", s.ctx, UserID).Return(nil, repoErr)

	movies, err := s.catalogRepository.GetRecommendedMovies(s.ctx, UserID, Limit)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(movies)
}

func (s *ServiceSuite) TestGetRecommendedMoviesGetMoviesError() {
	var (
		repoErr = gofakeit.Error()
		UserID  = gofakeit.UUID()
		Limit   = int32(10)

		userGenres = []model.Genre{
			{GenreID: gofakeit.IntRange(1, 20), GenreName: gofakeit.MovieGenre()},
		}
	)

	s.catalogRepository.On("GetUsersFavouriteGenres", s.ctx, UserID).Return(userGenres, nil)
	s.catalogRepository.On("GetRecommendedMovies", s.ctx, userGenres, Limit).Return(nil, repoErr)

	movies, err := s.catalogRepository.GetRecommendedMovies(s.ctx, UserID, Limit)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(movies)
}

func (s *ServiceSuite) TestGetRecommendedMoviesEmptyGenres() {
	var (
		UserID = gofakeit.UUID()
		Limit  = int32(10)

		emptyGenres = []model.Genre{}
	)

	s.catalogRepository.On("GetUsersFavouriteGenres", s.ctx, UserID).Return(emptyGenres, nil)
	s.catalogRepository.On("GetRecommendedMovies", s.ctx, emptyGenres, Limit).Return([]*model.Movie{}, nil)

	movies, err := s.catalogRepository.GetRecommendedMovies(s.ctx, UserID, Limit)
	s.Require().NoError(err)
	s.Require().Empty(movies)
}
