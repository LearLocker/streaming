package subscription

import (
	"github.com/LearLocker/streaming/subscription/internal/clients"
	"github.com/LearLocker/streaming/subscription/internal/repository"
	def "github.com/LearLocker/streaming/subscription/internal/service"
)

var _ def.SubscriptionService = (*Service)(nil)

type Service struct {
	subscriptionRepository repository.SubscriptionRepository
	catalogClient          clients.CatalogClient // интерфейс, а не *clients.catalogClient
}

func NewService(
	subscriptionRepository repository.SubscriptionRepository,
	catalogClient clients.CatalogClient,
) *Service {
	return &Service{
		subscriptionRepository: subscriptionRepository,
		catalogClient:          catalogClient,
	}
}
