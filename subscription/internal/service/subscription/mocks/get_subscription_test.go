package subscription

import (
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceSuite) TestGetSubscriptionSuccess() {
	var (
		UUID              = gofakeit.UUID()
		modelSubscription = model.CreateSubscriptionInfo{
			PlanID:        UUID,
			PaymentMethod: RandomPaymentMethod(),
		}
	)

	s.subscriptionRepository.On("GetByUUID", s.ctx, UUID).Return(modelSubscription, nil)

	res, err := s.subscriptionRepository.GetByUUID(s.ctx, UUID)
	s.Require().Error(err)
	s.Require().Nil(res)
}

func (s *ServiceSuite) TestGetSubscriptionNotFound() {
	var (
		UUID = gofakeit.UUID()
	)

	s.subscriptionRepository.On("GetByUUID", s.ctx, UUID).Return(nil, model.ErrSubscriptionNotFound)

	res, err := s.subscriptionRepository.GetByUUID(s.ctx, UUID)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *ServiceSuite) TestGetSubscriptionServiceError() {
	var (
		repoErr = gofakeit.Error()
		UUID    = gofakeit.UUID()
	)

	s.subscriptionRepository.On("GetByUUID", s.ctx, UUID).Return(nil, repoErr)

	res, err := s.subscriptionRepository.GetByUUID(s.ctx, UUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(res)
}
