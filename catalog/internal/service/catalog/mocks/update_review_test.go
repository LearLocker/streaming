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

	s.catalogRepository.On("UpdateReview", s.ctx, imdbId, authorId, text, ranking).Return(nil)

	res, err := s.service.UpdateReview(s.ctx, imdbId, authorId, text)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(updateReviewModel, res)
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

	s.catalogRepository.On("UpdateReview", s.ctx, imdbId, authorId, text, ranking).Return(model.ErrReviewNotFound)

	res, err := s.service.UpdateReview(s.ctx, imdbId, authorId, text)
	s.Require().Error(err)
	s.Require().Nil(res)

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

	s.catalogRepository.On("UpdateReview", s.ctx, imdbId, authorId, text, ranking).Return(repoErr)

	res, err := s.service.UpdateReview(s.ctx, imdbId, authorId, text)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(res)
}
