package subscription

import (
	"sync"

	"github.com/LearLocker/streaming/subscription/internal/model"
	def "github.com/LearLocker/streaming/subscription/internal/repository"
)

var _ def.SubscriptionRepository = (*repository)(nil)

type repository struct {
	mu            sync.RWMutex
	subscriptions map[string]*model.Subscription
}

func NewRepository() *repository {
	return &repository{
		subscriptions: make(map[string]*model.Subscription),
	}
}
