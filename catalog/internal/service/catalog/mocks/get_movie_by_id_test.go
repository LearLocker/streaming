package catalog

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceSuite) TestGetMovieByIdSuccess() {
	var (
		imdbId       = gofakeit.UUID()
		title        = gofakeit.MovieName()
		posterPath   = gofakeit.ImageURL(1920, 1080)
		youtubeID    = "dQw4w9WgXcQ"
		genreName    = gofakeit.MovieGenre()
		rankingName  = gofakeit.Paragraph(3, 5, 5, " ")
		rankingValue = gofakeit.IntRange(1, 100)

		modelMovie = model.Movie{
			ImdbID:     imdbId,
			Title:      title,
			PosterPath: posterPath,
			YouTubeID:  youtubeID,
			Genre: []model.Genre{{
				GenreID:   gofakeit.IntRange(1, 100),
				GenreName: genreName,
			}},
			Ranking: model.Ranking{
				RankingValue: rankingValue,
				RankingName:  rankingName,
			},
		}
	)

	s.catalogRepository.On("GetMovieById", s.ctx, imdbId).Return(modelMovie, nil)

	res, err := s.service.GetMovieById(s.ctx, imdbId)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(modelMovie, res)
}

func (s *ServiceSuite) TestGetMovieByIdNotFound() {
	var (
		imdbId = gofakeit.UUID()
	)

	s.catalogRepository.On("GetMovieById", s.ctx, imdbId).Return(nil, model.ErrMovieNotFound)

	res, err := s.service.GetMovieById(s.ctx, imdbId)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *ServiceSuite) TestGetMovieByIdRepoError() {
	var (
		imdbId  = gofakeit.UUID()
		repoErr = gofakeit.Error()
	)

	s.catalogRepository.On("GetMovieById", s.ctx, imdbId).Return(nil, repoErr)

	res, err := s.service.GetMovieById(s.ctx, imdbId)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(res)
}
