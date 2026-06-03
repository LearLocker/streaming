package subscription

import (
	"context"
	"testing"

	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/LearLocker/streaming/subscription/internal/repository/mocks"
	subscriptionService "github.com/LearLocker/streaming/subscription/internal/service/subscription"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	subscriptionRepository *mocks.SubscriptionRepository

	service *subscriptionService.Service

	catalogClient *mocks.CatalogClient
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.subscriptionRepository = mocks.NewSubscriptionRepository(s.T())
	s.catalogClient = mocks.NewCatalogClient(s.T())

	s.service = subscriptionService.NewService(
		s.subscriptionRepository,
		s.catalogClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func RandomPaymentMethod() model.PaymentMethod {
	methods := []model.PaymentMethod{
		model.PaymentMethodCard,
		model.PaymentMethodSBP,
		model.PaymentMethodWallet,
	}
	return methods[gofakeit.IntRange(0, len(methods)-1)]
}
