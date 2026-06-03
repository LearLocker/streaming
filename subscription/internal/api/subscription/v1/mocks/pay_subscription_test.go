package v1

import (
	subscriptionV1 "github.com/LearLocker/streaming/shared/pkg/openapi/subscription/v1"
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestPaySubscriptionSuccess() {
	var (
		planID = gofakeit.UUID()

		paySubscriptionRequest = &subscriptionV1.PaySubscriptionRequest{
			PaymentMethod: subscriptionV1.PaymentMethod(RandomPaymentMethod()),
		}

		paySubscriptionParams = subscriptionV1.PaySubscriptionByUuidParams{
			UUID: planID,
		}

		expectedModel = model.Subscription{
			UUID:          gofakeit.UUID(),
			PlanID:        planID,
			PlanName:      gofakeit.Word(),
			Amount:        int64(gofakeit.Price(12, 100)),
			Currency:      gofakeit.CurrencyShort(),
			PaymentMethod: RandomPaymentMethod(),
			Status:        model.StatusActive,
		}
	)

	s.subscriptionService.On("PaySubscriptionByUuid", s.ctx, paySubscriptionRequest, paySubscriptionParams).Return(expectedModel, nil)

	res, err := s.api.PaySubscriptionByUuid(s.ctx, paySubscriptionRequest, paySubscriptionParams)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedModel, res)
}

func (s *APISuite) TestPaySubscriptionNotFound() {
	var (
		paySubscriptionRequest = &subscriptionV1.PaySubscriptionRequest{
			PaymentMethod: subscriptionV1.PaymentMethod(RandomPaymentMethod()),
		}

		paySubscriptionParams = subscriptionV1.PaySubscriptionByUuidParams{
			UUID: gofakeit.UUID(),
		}
	)

	s.subscriptionService.On("PaySubscriptionByUuid", s.ctx, paySubscriptionRequest, paySubscriptionParams).Return(nil, model.ErrSubscriptionNotFound)

	res, err := s.api.PaySubscriptionByUuid(s.ctx, paySubscriptionRequest, paySubscriptionParams)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestPaySubscriptionServiceError() {
	var (
		serviceErr = gofakeit.Error()

		paySubscriptionRequest = &subscriptionV1.PaySubscriptionRequest{
			PaymentMethod: subscriptionV1.PaymentMethod(RandomPaymentMethod()),
		}

		paySubscriptionParams = subscriptionV1.PaySubscriptionByUuidParams{
			UUID: gofakeit.UUID(),
		}
	)

	s.subscriptionService.On("PaySubscriptionByUuid", s.ctx, paySubscriptionRequest, paySubscriptionParams).Return(nil, serviceErr)

	res, err := s.api.PaySubscriptionByUuid(s.ctx, paySubscriptionRequest, paySubscriptionParams)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
