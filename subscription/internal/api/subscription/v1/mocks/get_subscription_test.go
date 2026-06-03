package v1

import (
	subscriptionV1 "github.com/LearLocker/streaming/shared/pkg/openapi/subscription/v1"
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestGetSubscriptionSuccess() {
	var (
		getSubscriptionReq = subscriptionV1.GetSubscriptionByUuidParams{
			UUID: gofakeit.UUID(),
		}

		modelSubscription = model.CreateSubscriptionInfo{
			PlanID:        gofakeit.UUID(),
			PaymentMethod: RandomPaymentMethod(),
		}
	)

	s.subscriptionService.On("GetSubscriptionByUuid", s.ctx, getSubscriptionReq).Return(modelSubscription, nil)

	res, err := s.api.GetSubscriptionByUuid(s.ctx, getSubscriptionReq)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(modelSubscription, res)
}

func (s *APISuite) TestGetSubscriptionNotFound() {
	var (
		getSubscriptionReq = subscriptionV1.GetSubscriptionByUuidParams{
			UUID: gofakeit.UUID(),
		}
	)

	s.subscriptionService.On("GetSubscriptionByUuid", s.ctx, getSubscriptionReq).Return(nil, model.ErrSubscriptionNotFound)

	res, err := s.api.GetSubscriptionByUuid(s.ctx, getSubscriptionReq)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetSubscriptionServiceError() {
	var (
		serviceErr = gofakeit.Error()

		getSubscriptionReq = subscriptionV1.GetSubscriptionByUuidParams{
			UUID: gofakeit.UUID(),
		}
	)

	s.subscriptionService.On("GetSubscriptionByUuid", s.ctx, getSubscriptionReq).Return(nil, serviceErr)

	res, err := s.api.GetSubscriptionByUuid(s.ctx, getSubscriptionReq)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
