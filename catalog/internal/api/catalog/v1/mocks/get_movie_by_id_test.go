package v1

import (
	"github.com/LearLocker/streaming/catalog/internal/converter"
	"github.com/LearLocker/streaming/catalog/internal/model"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestGetMovieByIdSuccess() {
	var (
		imdbId       = gofakeit.UUID()
		title        = gofakeit.MovieName()
		posterPath   = gofakeit.ImageURL(1920, 1080)
		youtubeID    = "dQw4w9WgXcQ"
		genreName    = gofakeit.MovieGenre()
		rankingName  = gofakeit.Paragraph(3, 5, 5, " ")
		rankingValue = gofakeit.IntRange(1, 100)

		req = &catalogV1.GetMovieByIdRequest{
			ImdbId: imdbId,
		}

		modelMovie = &model.Movie{
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
		expectedProtoMovie = converter.MovieToProto(modelMovie)
	)

	s.catalogService.On("GetMovieById", s.ctx, imdbId).Return(modelMovie, nil)

	res, err := s.api.GetMovieById(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedProtoMovie.ImdbId, res.GetMovie().GetImdbId())
}

func (s *APISuite) TestGetMovieByIdNotFound() {
	var (
		imdbId = gofakeit.UUID()

		req = &catalogV1.GetMovieByIdRequest{
			ImdbId: imdbId,
		}
	)

	s.catalogService.On("GetMovieById", s.ctx, imdbId).Return(model.Movie{}, model.ErrMovieNotFound)

	res, err := s.api.GetMovieById(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetMovieByIdServiceError() {
	var (
		imdbId     = gofakeit.UUID()
		serviceErr = gofakeit.Error()

		req = &catalogV1.GetMovieByIdRequest{
			ImdbId: imdbId,
		}
	)

	s.catalogService.On("GetMovieById", s.ctx, imdbId).Return(model.Movie{}, serviceErr)

	res, err := s.api.GetMovieById(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
