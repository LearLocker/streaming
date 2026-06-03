package catalog

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceSuite) TestUpdateReviewSuccess() {
	var (
		imdbId       = gofakeit.UUID()
		text         = gofakeit.Paragraph(3, 5, 5, " ")
		authorId     = gofakeit.BookAuthor()
		rankingName  = gofakeit.Word()
		rankingValue = int32(gofakeit.IntRange(1, 100))

		ranking = model.Ranking{
			RankingName:  rankingName,
			RankingValue: int(rankingValue),
		}

		updateReviewModel = model.UpdateReviewResult{
			RankingName:  rankingName,
			RankingValue: rankingValue,
			Text:         text,
		}
	)

	s.catalogRepository.On("UpdateReview", s.ctx, imdbId, authorId, text, ranking).Return(updateReviewModel, nil)

	err := s.catalogRepository.UpdateReview(s.ctx, imdbId, authorId, text, ranking)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestUpdateReviewNotFound() {
	var (
		imdbId   = gofakeit.UUID()
		text     = gofakeit.Paragraph(3, 5, 5, " ")
		authorId = gofakeit.BookAuthor()

		ranking = model.Ranking{
			RankingName:  gofakeit.Word(),
			RankingValue: gofakeit.IntRange(1, 100),
		}
	)

	s.catalogRepository.On("UpdateReview", s.ctx, imdbId, authorId, text, ranking).Return(nil, model.ErrReviewNotFound)

	err := s.catalogRepository.UpdateReview(s.ctx, imdbId, authorId, text, ranking)
	s.Require().Error(err)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *ServiceSuite) TestUpdateReviewRepoError() {
	var (
		repoErr  = gofakeit.Error()
		imdbId   = gofakeit.UUID()
		text     = gofakeit.Paragraph(3, 5, 5, " ")
		authorId = gofakeit.BookAuthor()

		ranking = model.Ranking{
			RankingName:  gofakeit.Word(),
			RankingValue: gofakeit.IntRange(1, 100),
		}
	)

	s.catalogRepository.On("UpdateReview", s.ctx, imdbId, authorId, text, ranking).Return(nil, repoErr)

	err := s.catalogRepository.UpdateReview(s.ctx, imdbId, authorId, text, ranking)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
