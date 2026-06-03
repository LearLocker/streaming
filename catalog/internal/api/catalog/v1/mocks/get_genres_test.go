package v1

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *APISuite) TestGetGenresSuccess() {
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

		req = &catalogV1.GetGenresRequest{}
	)

	s.catalogService.On("GetGenres", s.ctx, req).Return(expectedGenres, nil)

	res, err := s.api.GetGenres(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Len(res.GetGenres(), len(expectedGenres))
}

func (s *APISuite) TestGetGenresServiceError() {
	var (
		serviceErr = gofakeit.Error()
		req        = &catalogV1.GetGenresRequest{}
	)

	s.catalogService.On("GetGenres", s.ctx, req).Return(nil, serviceErr)

	res, err := s.api.GetGenres(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}

func (s *APISuite) TestGetGenresEmptyList() {
	var (
		expectedGenres = []model.Genre{}
		req            = &catalogV1.GetGenresRequest{}
	)

	s.catalogService.On("GetGenres", s.ctx, req).Return(expectedGenres, nil)

	res, err := s.api.GetGenres(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Empty(res.GetGenres())
}
