package catalog

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestGetGenresSuccess() {
	var (
		expectedGenres = []model.Genre{
			{
				GenreID:   gofakeit.IntRange(1, 100),
				GenreName: gofakeit.MovieGenre(),
			},
			{
				GenreID:   gofakeit.IntRange(1, 100),
				GenreName: gofakeit.MovieGenre(),
			},
		}
	)

	s.catalogRepository.On("GetGenres", s.ctx).Return(expectedGenres, nil)

	genres, err := s.service.GetGenres(s.ctx)
	s.Require().NoError(err)
	s.Require().NotNil(genres)
	s.Require().Len(genres, len(expectedGenres))
	s.Require().Equal(expectedGenres, genres)
}

func (s *ServiceSuite) TestGetGenresRepoError() {
	var (
		repoErr = gofakeit.Error()
	)

	s.catalogRepository.On("GetGenres", s.ctx).Return(nil, repoErr)

	genres, err := s.service.GetGenres(s.ctx)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(genres)
}

func (s *ServiceSuite) TestGetGenresEmptyList() {
	s.catalogRepository.On("GetGenres", s.ctx).Return([]model.Genre{}, nil)

	genres, err := s.service.GetGenres(s.ctx)
	s.Require().NoError(err)
	s.Require().NotNil(genres)
	s.Require().Empty(genres)
}
