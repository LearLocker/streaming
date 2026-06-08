package subscription

import (
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceSuite) TestCreateSubscriptionSuccess() {
	var (
		planID = gofakeit.UUID()

		modelSubscription = &model.Subscription{
			UUID:          gofakeit.UUID(),
			PlanID:        planID,
			PlanName:      gofakeit.Word(),
			Amount:        int64(gofakeit.Price(12, 100)),
			Currency:      gofakeit.CurrencyShort(),
			PaymentMethod: RandomPaymentMethod(),
			Status:        model.StatusPending,
		}

		info = model.CreateSubscriptionInfo{
			PlanID:        planID,
			PaymentMethod: RandomPaymentMethod(),
		}
	)

	s.subscriptionRepository.On("Create", s.ctx, modelSubscription).Return(nil)

	res, err := s.service.CreateSubscription(s.ctx, info)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(modelSubscription.PlanID, res.PlanID)
	s.Require().Equal(modelSubscription.PlanName, res.PlanName)
	s.Require().Equal(modelSubscription.Amount, res.Amount)
	s.Require().Equal(modelSubscription.Currency, res.Currency)
	s.Require().Equal(info.PaymentMethod, res.PaymentMethod)
	s.Require().Equal(model.StatusPending, res.Status)
	s.Require().NotEmpty(res.UUID)
}

func (s *ServiceSuite) TestCreateSubscriptionNotFound() {
	var (
		planID = gofakeit.UUID()

		modelSubscription = &model.Subscription{
			UUID:          gofakeit.UUID(),
			PlanID:        planID,
			PlanName:      gofakeit.Word(),
			Amount:        int64(gofakeit.Price(12, 100)),
			Currency:      gofakeit.CurrencyShort(),
			PaymentMethod: RandomPaymentMethod(),
			Status:        model.StatusPending,
		}

		info = model.CreateSubscriptionInfo{
			PlanID:        planID,
			PaymentMethod: RandomPaymentMethod(),
		}
	)

	s.subscriptionRepository.On("Create", s.ctx, modelSubscription).Return(nil, model.ErrSubscriptionNotFound)

	res, err := s.service.CreateSubscription(s.ctx, info)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *ServiceSuite) TestCreateSubscriptionServiceError() {
	var (
		repoErr = gofakeit.Error()
		planID  = gofakeit.UUID()

		modelSubscription = &model.Subscription{
			UUID:          gofakeit.UUID(),
			PlanID:        planID,
			PlanName:      gofakeit.Word(),
			Amount:        int64(gofakeit.Price(12, 100)),
			Currency:      gofakeit.CurrencyShort(),
			PaymentMethod: RandomPaymentMethod(),
			Status:        model.StatusPending,
		}

		info = model.CreateSubscriptionInfo{
			PlanID:        planID,
			PaymentMethod: RandomPaymentMethod(),
		}
	)

	s.subscriptionRepository.On("Create", s.ctx, modelSubscription).Return(nil, repoErr)

	res, err := s.service.CreateSubscription(s.ctx, info)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(res)
}
