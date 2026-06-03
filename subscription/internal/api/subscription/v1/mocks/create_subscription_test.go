package v1

import (
	subscriptionV1 "github.com/LearLocker/streaming/shared/pkg/openapi/subscription/v1"
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestCreateSubscriptionSuccess() {
	var (
		planID = gofakeit.UUID()

		createSubscriptionReq = &subscriptionV1.CreateSubscriptionRequest{
			PlanID:        uuid.MustParse(planID),
			PaymentMethod: subscriptionV1.PaymentMethod(RandomPaymentMethod()),
		}

		expectedModel = model.Subscription{
			UUID:          gofakeit.UUID(),
			PlanID:        planID,
			PlanName:      gofakeit.Word(),
			Amount:        int64(gofakeit.Price(12, 100)),
			Currency:      gofakeit.CurrencyShort(),
			PaymentMethod: RandomPaymentMethod(),
			Status:        model.StatusPending,
		}
	)

	s.subscriptionService.On("CreateSubscription", s.ctx, createSubscriptionReq).Return(expectedModel, nil)

	res, err := s.api.CreateSubscription(s.ctx, createSubscriptionReq)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedModel, res)
}

func (s *APISuite) TestCreateSubscriptionNotFound() {
	var (
		createSubscriptionReq = &subscriptionV1.CreateSubscriptionRequest{
			PlanID:        uuid.MustParse(gofakeit.UUID()),
			PaymentMethod: subscriptionV1.PaymentMethod(RandomPaymentMethod()),
		}
	)

	s.subscriptionService.On("CreateSubscription", s.ctx, createSubscriptionReq).Return(nil, model.ErrSubscriptionNotFound)

	res, err := s.api.CreateSubscription(s.ctx, createSubscriptionReq)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestCreateSubscriptionServiceError() {
	var (
		serviceErr = gofakeit.Error()

		createSubscriptionReq = &subscriptionV1.CreateSubscriptionRequest{
			PlanID:        uuid.MustParse(gofakeit.UUID()),
			PaymentMethod: subscriptionV1.PaymentMethod(RandomPaymentMethod()),
		}
	)

	s.subscriptionService.On("CreateSubscription", s.ctx, createSubscriptionReq).Return(nil, serviceErr)

	res, err := s.api.CreateSubscription(s.ctx, createSubscriptionReq)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
