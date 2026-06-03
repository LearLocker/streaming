package v1

import (
	"github.com/LearLocker/streaming/catalog/internal/converter"
	"github.com/brianvoe/gofakeit/v6"

	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
)

func (s *APISuite) TestAddMovieSuccess() {
	var (
		imdbId       = gofakeit.UUID()
		title        = gofakeit.MovieName()
		posterPath   = gofakeit.ImageURL(1920, 1080)
		youtubeID    = "dQw4w9WgXcQ"
		genreName    = gofakeit.MovieGenre()
		rankingName  = gofakeit.Paragraph(3, 5, 5, " ")
		rankingValue = int32(gofakeit.IntRange(1, 100))

		expectedUUID = gofakeit.UUID()

		req = &catalogV1.AddMovieRequest{
			ImdbId:     imdbId,
			Title:      title,
			PosterPath: posterPath,
			YoutubeId:  youtubeID,
			Genre: []*catalogV1.Genre{{
				GenreId:   int32(gofakeit.IntRange(1, 100)),
				GenreName: genreName,
			}},
			Ranking: &catalogV1.Ranking{
				RankingValue: rankingValue,
				RankingName:  rankingName,
			},
		}

		expectedModelInfo = converter.AddMovieInfoFromProto(req)
	)

	s.catalogService.On("AddMovie", s.ctx, expectedModelInfo).Return(expectedUUID, nil)

	res, err := s.api.AddMovie(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedUUID, res.GetInsertedId())
}

func (s *APISuite) TestAddMovieServiceError() {
	var (
		serviceErr   = gofakeit.Error()
		imdbId       = gofakeit.UUID()
		title        = gofakeit.MovieName()
		posterPath   = gofakeit.ImageURL(1920, 1080)
		youtubeID    = "dQw4w9WgXcQ"
		genreName    = gofakeit.MovieGenre()
		rankingName  = gofakeit.Paragraph(3, 5, 5, " ")
		rankingValue = int32(gofakeit.IntRange(1, 100))

		req = &catalogV1.AddMovieRequest{
			ImdbId:     imdbId,
			Title:      title,
			PosterPath: posterPath,
			YoutubeId:  youtubeID,
			Genre: []*catalogV1.Genre{{
				GenreId:   int32(gofakeit.IntRange(1, 100)),
				GenreName: genreName,
			}},
			Ranking: &catalogV1.Ranking{
				RankingValue: rankingValue,
				RankingName:  rankingName,
			},
		}

		expectedModelInfo = converter.AddMovieInfoFromProto(req)
	)

	s.catalogService.On("AddMovie", s.ctx, expectedModelInfo).Return("", serviceErr)

	res, err := s.api.AddMovie(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
