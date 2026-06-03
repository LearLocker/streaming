package subscription

import (
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/brianvoe/gofakeit/v6"
)

func (s *ServiceSuite) TestPaySubscriptionSuccess() {
	var (
		UUID   = gofakeit.UUID()
		Status = model.StatusActive
	)

	s.subscriptionRepository.On("UpdateStatus", s.ctx, UUID, Status, nil).Return(nil)

	err := s.subscriptionRepository.UpdateStatus(s.ctx, UUID, Status, nil)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestPaySubscriptionNotFound() {
	var (
		UUID   = gofakeit.UUID()
		Status = model.StatusActive
	)

	s.subscriptionRepository.On("UpdateStatus", s.ctx, UUID, Status, nil).Return(model.ErrSubscriptionNotFound)

	err := s.subscriptionRepository.UpdateStatus(s.ctx, UUID, Status, nil)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrSubscriptionNotFound)
}

func (s *ServiceSuite) TestPaySubscriptionServiceError() {
	var (
		repoErr = gofakeit.Error()
		UUID    = gofakeit.UUID()
		Status  = model.StatusActive
	)

	s.subscriptionRepository.On("UpdateStatus", s.ctx, UUID, Status, nil).Return(repoErr)

	err := s.subscriptionRepository.UpdateStatus(s.ctx, UUID, Status, nil)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
