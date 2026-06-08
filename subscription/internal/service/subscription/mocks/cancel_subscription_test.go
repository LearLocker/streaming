package subscription

import (
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestCancelSubscriptionSuccess() {
	var (
		UUID   = gofakeit.UUID()
		Status = model.StatusCancelled
	)

	s.subscriptionRepository.On("UpdateStatus", s.ctx, UUID, Status, nil).Return(nil)

	err := s.service.CancelSubscription(s.ctx, UUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelSubscriptionNotFound() {
	var (
		UUID   = gofakeit.UUID()
		Status = model.StatusCancelled
	)

	s.subscriptionRepository.On("UpdateStatus", s.ctx, UUID, Status, nil).Return(model.ErrSubscriptionNotFound)

	err := s.service.CancelSubscription(s.ctx, UUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrSubscriptionNotFound)
}

func (s *ServiceSuite) TestCancelSubscriptionServiceError() {
	var (
		repoErr = gofakeit.Error()
		UUID    = gofakeit.UUID()
		Status  = model.StatusCancelled
	)

	s.subscriptionRepository.On("UpdateStatus", s.ctx, UUID, Status, nil).Return(repoErr)

	err := s.service.CancelSubscription(s.ctx, UUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
