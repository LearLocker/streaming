package v1

import (
	"github.com/LearLocker/streaming/payment/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceSuite) TestGetPaymentSuccess() {
	var (
		paymentID = gofakeit.UUID()

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
	s.paymentRepository.On("GetByID", s.ctx, paymentID).Return(modelPayment, nil)

	res, err := s.paymentRepository.GetByID(s.ctx, paymentID)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(modelPayment, res)
}

func (s *ServiceSuite) TestGetPaymentNotFound() {
	var (
		paymentID = gofakeit.UUID()
	)

	s.paymentRepository.On("GetByID", s.ctx, paymentID).Return(nil, model.ErrPaymentNotFound)

	res, err := s.paymentRepository.GetByID(s.ctx, paymentID)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *ServiceSuite) TestGetPaymentRepoError() {
	var (
		paymentID = gofakeit.UUID()
		repoErr   = gofakeit.Error()
	)

	s.paymentRepository.On("GetByID", s.ctx, paymentID).Return(model.Payment{}, repoErr)

	res, err := s.paymentRepository.GetByID(s.ctx, paymentID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(res)
}
