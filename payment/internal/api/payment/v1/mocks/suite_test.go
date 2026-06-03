package v1

import (
	"context"
	"testing"

	apiV1 "github.com/LearLocker/streaming/payment/internal/api/payment/v1"
	"github.com/LearLocker/streaming/payment/internal/model"
	"github.com/LearLocker/streaming/payment/internal/service/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	paymentService *mocks.PaymentService

	api *apiV1.Api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.paymentService = mocks.NewPaymentService(s.T())

	s.api = apiV1.NewAPI(
		s.paymentService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
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
