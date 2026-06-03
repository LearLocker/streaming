package payment

import (
	"github.com/LearLocker/streaming/payment/internal/repository"
	def "github.com/LearLocker/streaming/payment/internal/service"
)

var _ def.PaymentService = (*Service)(nil)

type Service struct {
	paymentRepository repository.PaymentRepository
}

func NewService(paymentRepository repository.PaymentRepository) *Service {
	return &Service{
		paymentRepository: paymentRepository,
	}
}
