package payment

import (
	"sync"

	"github.com/LearLocker/streaming/payment/internal/model"
	def "github.com/LearLocker/streaming/payment/internal/repository"
)

var _ def.PaymentRepository = (*repository)(nil)

type repository struct {
	mu       sync.RWMutex
	payments map[string]*model.Payment
}

func NewRepository() *repository {
	return &repository{
		payments: make(map[string]*model.Payment),
	}
}
