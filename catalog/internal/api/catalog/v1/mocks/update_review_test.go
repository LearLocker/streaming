package v1

import (
	"github.com/LearLocker/streaming/catalog/internal/model"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestUpdateReviewSuccess() {
	var (
		imdbId       = gofakeit.UUID()
		text         = gofakeit.Paragraph(3, 5, 5, " ")
		authorId     = gofakeit.BookAuthor()
		rankingName  = gofakeit.Word()
		rankingValue = int32(gofakeit.IntRange(1, 100))

		updateReviewModel = &model.UpdateReviewResult{
			RankingName:  rankingName,
			RankingValue: rankingValue,
			Text:         text,
		}

		req = &catalogV1.UpdateReviewRequest{
			ImdbId:   imdbId,
			AuthorId: text,
			Text:     text,
		}
	)

	s.catalogService.On("UpdateReview", s.ctx, imdbId, authorId, text).Return(updateReviewModel, nil)

	res, err := s.api.UpdateReview(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func (s *APISuite) TestUpdateReviewNilInfo() {
	res, err := s.api.UpdateReview(s.ctx, nil)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
}

func (s *APISuite) TestUpdateReviewNotFound() {
	var (
		imdbId   = gofakeit.UUID()
		text     = gofakeit.Paragraph(3, 5, 5, " ")
		authorId = gofakeit.BookAuthor()

		req = &catalogV1.UpdateReviewRequest{
			ImdbId:   imdbId,
			AuthorId: text,
			Text:     text,
		}
	)

	s.catalogService.On("UpdateReview", s.ctx, imdbId, authorId, text).Return(nil, model.ErrReviewNotFound)

	res, err := s.api.UpdateReview(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestUpdateReviewServiceError() {
	var (
		serviceErr = gofakeit.Error()
		imdbId     = gofakeit.UUID()
		text       = gofakeit.Paragraph(3, 5, 5, " ")
		authorId   = gofakeit.BookAuthor()

		req = &catalogV1.UpdateReviewRequest{
			ImdbId:   imdbId,
			AuthorId: text,
			Text:     text,
		}
	)

	s.catalogService.On("UpdateReview", s.ctx, imdbId, authorId, text).Return(nil, serviceErr)

	res, err := s.api.UpdateReview(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
