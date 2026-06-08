package v1

import (
	"github.com/LearLocker/streaming/payment/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceSuite) TestProcessPaymentSuccess() {
	var (
		paymentID = gofakeit.UUID()

		processPaymentInfo = model.ProcessPaymentInfo{
			SubscriptionID: gofakeit.UUID(),
			UserID:         gofakeit.UUID(),
			Amount:         int64(gofakeit.Price(12, 100)),
			Currency:       gofakeit.CurrencyShort(),
			Method:         gofakeit.Word(),
		}

		modelPayment = model.Payment{
			PaymentID:      paymentID,
			SubscriptionID: gofakeit.UUID(),
			UserID:         gofakeit.UUID(),
			Amount:         int64(gofakeit.Price(12, 100)),
			Currency:       gofakeit.CurrencyShort(),
			Method:         gofakeit.Word(),
			Status:         RandomPaymentStatus(),
			Reason:         gofakeit.Word(),
		}
	)

	s.paymentRepository.On("Create", s.ctx, processPaymentInfo).Return(modelPayment, nil)

	res, err := s.service.ProcessPayment(s.ctx, processPaymentInfo)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(modelPayment, res)
}

func (s *ServiceSuite) TestProcessPaymentNotFound() {
	var (
		processPaymentInfo = model.ProcessPaymentInfo{
			SubscriptionID: gofakeit.UUID(),
			UserID:         gofakeit.UUID(),
			Amount:         int64(gofakeit.Price(12, 100)),
			Currency:       gofakeit.CurrencyShort(),
			Method:         gofakeit.Word(),
		}
	)

	s.paymentRepository.On("Create", s.ctx, processPaymentInfo).Return(model.Payment{}, nil)

	res, err := s.service.ProcessPayment(s.ctx, processPaymentInfo)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *ServiceSuite) TestProcessPaymentRepoError() {
	var (
		repoErr = gofakeit.Error()

		processPaymentInfo = model.ProcessPaymentInfo{
			SubscriptionID: gofakeit.UUID(),
			UserID:         gofakeit.UUID(),
			Amount:         int64(gofakeit.Price(12, 100)),
			Currency:       gofakeit.CurrencyShort(),
			Method:         gofakeit.Word(),
		}
	)

	s.paymentRepository.On("Create", s.ctx, processPaymentInfo).Return(nil, repoErr)

	res, err := s.service.ProcessPayment(s.ctx, processPaymentInfo)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(res)
}
