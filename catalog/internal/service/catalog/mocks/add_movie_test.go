package catalog

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestAddMovieSuccess() {
	var (
		imdbId       = gofakeit.UUID()
		title        = gofakeit.MovieName()
		posterPath   = gofakeit.ImageURL(1920, 1080)
		youtubeID    = "dQw4w9WgXcQ"
		genreName    = gofakeit.MovieGenre()
		rankingName  = gofakeit.Paragraph(3, 5, 5, " ")
		rankingValue = gofakeit.IntRange(1, 100)

		expectedUUID = gofakeit.UUID()

		addMovie = model.AddMovie{
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

	s.catalogRepository.On("AddMovie", s.ctx, addMovie).Return(expectedUUID, nil)

	insertedId, err := s.catalogRepository.AddMovie(s.ctx, addMovie)
	s.Require().NoError(err)
	s.Require().NotNil(insertedId)
	s.Require().Equal(expectedUUID, insertedId)
}

func (s *ServiceSuite) TestAddMovieRepoError() {
	var (
		repoErr      = gofakeit.Error()
		imdbId       = gofakeit.UUID()
		title        = gofakeit.MovieName()
		posterPath   = gofakeit.ImageURL(1920, 1080)
		youtubeID    = "dQw4w9WgXcQ"
		genreName    = gofakeit.MovieGenre()
		rankingName  = gofakeit.Paragraph(3, 5, 5, " ")
		rankingValue = gofakeit.IntRange(1, 100)

		addMovie = model.AddMovie{
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

	s.catalogRepository.On("AddMovie", s.ctx, addMovie).Return("", repoErr)

	res, err := s.catalogRepository.AddMovie(s.ctx, addMovie)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(res)
}
