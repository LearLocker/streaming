package v1

import (
	"context"
	"testing"

	apiV1 "github.com/LearLocker/streaming/subscription/internal/api/subscription/v1"
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/LearLocker/streaming/subscription/internal/service/mocks"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	subscriptionService *mocks.SubscriptionService

	api *apiV1.Api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.subscriptionService = mocks.NewSubscriptionService(s.T())

	s.api = apiV1.NewAPI(
		s.subscriptionService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}

func RandomPaymentMethod() model.PaymentMethod {
	methods := []model.PaymentMethod{
		model.PaymentMethodCard,
		model.PaymentMethodSBP,
		model.PaymentMethodWallet,
	}
	return methods[gofakeit.IntRange(0, len(methods)-1)]
}
