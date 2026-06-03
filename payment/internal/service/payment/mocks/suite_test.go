package v1

import (
	"context"
	"testing"

	"github.com/LearLocker/streaming/payment/internal/model"
	"github.com/LearLocker/streaming/payment/internal/repository/mocks"
	paymentService "github.com/LearLocker/streaming/payment/internal/service/payment"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	paymentRepository *mocks.PaymentRepository

	service *paymentService.Service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.paymentRepository = mocks.NewPaymentRepository(s.T())

	s.service = paymentService.NewService(
		s.paymentRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func RandomPaymentStatus() model.PaymentStatus {
	statuses := []model.PaymentStatus{
		model.StatusPending,
		model.StatusProcessing,
		model.StatusSuccess,
		model.StatusFailed,
		model.StatusRefunded,
	}
	return statuses[gofakeit.IntRange(0, len(statuses)-1)]
}
