package v1

import (
	subscriptionV1 "github.com/LearLocker/streaming/shared/pkg/openapi/subscription/v1"
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestCancelSubscriptionSuccess() {
	var (
		cancelSubscriptionReq = subscriptionV1.CancelSubscriptionByUuidParams{
			UUID: gofakeit.UUID(),
		}

		expectedResponse = &subscriptionV1.CancelSubscriptionByUuidNoContent{}
	)

	s.subscriptionService.On("CancelSubscriptionByUuid", s.ctx, cancelSubscriptionReq).Return(expectedResponse, nil)

	res, err := s.api.CancelSubscriptionByUuid(s.ctx, cancelSubscriptionReq)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedResponse, res)
}

func (s *APISuite) TestCancelSubscriptionNotFound() {
	var (
		cancelSubscriptionReq = subscriptionV1.CancelSubscriptionByUuidParams{
			UUID: gofakeit.UUID(),
		}
	)

	s.subscriptionService.On("CancelSubscriptionByUuid", s.ctx, cancelSubscriptionReq).Return(nil, model.ErrSubscriptionNotFound)

	res, err := s.api.CancelSubscriptionByUuid(s.ctx, cancelSubscriptionReq)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestCancelSubscriptionServiceError() {
	var (
		serviceErr = gofakeit.Error()

		cancelSubscriptionReq = subscriptionV1.CancelSubscriptionByUuidParams{
			UUID: gofakeit.UUID(),
		}
	)

	s.subscriptionService.On("CancelSubscriptionByUuid", s.ctx, cancelSubscriptionReq).Return(nil, serviceErr)

	res, err := s.api.CancelSubscriptionByUuid(s.ctx, cancelSubscriptionReq)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
