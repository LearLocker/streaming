package subscription

import (
	"github.com/LearLocker/streaming/subscription/internal/clients"
	"github.com/LearLocker/streaming/subscription/internal/repository"
	def "github.com/LearLocker/streaming/subscription/internal/service"
)

var _ def.SubscriptionService = (*service)(nil)

type service struct {
	subscriptionRepository repository.SubscriptionRepository
	catalogClient          *clients.CatalogClient
}

func NewService(
	subscriptionRepository repository.SubscriptionRepository,
	catalogClient *clients.CatalogClient,
) *service {
	return &service{
		subscriptionRepository: subscriptionRepository,
		catalogClient:          catalogClient,
	}
}
