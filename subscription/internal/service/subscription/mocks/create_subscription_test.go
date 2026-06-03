package subscription

import (
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestCreateSubscriptionSuccess() {
	var (
		modelSubscription = &model.Subscription{
			UUID:          gofakeit.UUID(),
			PlanID:        gofakeit.UUID(),
			PlanName:      gofakeit.Word(),
			Amount:        int64(gofakeit.Price(12, 100)),
			Currency:      gofakeit.CurrencyShort(),
			PaymentMethod: RandomPaymentMethod(),
			Status:        model.StatusPending,
		}
	)

	s.subscriptionRepository.On("Create", s.ctx, modelSubscription).Return(nil)

	err := s.subscriptionRepository.Create(s.ctx, modelSubscription)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCreateSubscriptionNotFound() {
	var (
		modelSubscription = &model.Subscription{
			UUID:          gofakeit.UUID(),
			PlanID:        gofakeit.UUID(),
			PlanName:      gofakeit.Word(),
			Amount:        int64(gofakeit.Price(12, 100)),
			Currency:      gofakeit.CurrencyShort(),
			PaymentMethod: RandomPaymentMethod(),
			Status:        model.StatusPending,
		}
	)

	s.subscriptionRepository.On("Create", s.ctx, modelSubscription).Return(nil, model.ErrSubscriptionNotFound)

	err := s.subscriptionRepository.Create(s.ctx, modelSubscription)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrSubscriptionNotFound)
}

func (s *ServiceSuite) TestCreateSubscriptionServiceError() {
	var (
		repoErr           = gofakeit.Error()
		modelSubscription = &model.Subscription{
			UUID:          gofakeit.UUID(),
			PlanID:        gofakeit.UUID(),
			PlanName:      gofakeit.Word(),
			Amount:        int64(gofakeit.Price(12, 100)),
			Currency:      gofakeit.CurrencyShort(),
			PaymentMethod: RandomPaymentMethod(),
			Status:        model.StatusPending,
		}
	)

	s.subscriptionRepository.On("Create", s.ctx, modelSubscription).Return(nil, repoErr)

	err := s.subscriptionRepository.Create(s.ctx, modelSubscription)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
