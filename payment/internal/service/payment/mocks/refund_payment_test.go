package v1

import (
	"github.com/LearLocker/streaming/payment/internal/model"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServiceSuite) TestRefundPaymentSuccess() {
	var (
		paymentID = gofakeit.UUID()
		reason    = gofakeit.Word()

		refundPaymentInfo = model.RefundPaymentInfo{
			PaymentID: paymentID,
			Reason:    reason,
		}

		modelRefundPayment = model.Payment{
			PaymentID:      paymentID,
			SubscriptionID: gofakeit.UUID(),
			UserID:         gofakeit.UUID(),
			Amount:         int64(gofakeit.Price(12, 100)),
			Currency:       gofakeit.CurrencyShort(),
			Method:         gofakeit.Word(),
			Status:         RandomPaymentStatus(),
			Reason:         reason,
		}
	)

	s.paymentRepository.On("Refund", s.ctx, refundPaymentInfo).Return(modelRefundPayment, nil)

	res, err := s.service.RefundPayment(s.ctx, refundPaymentInfo)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(modelRefundPayment, res)
}

func (s *ServiceSuite) TestRefundPaymentNotFound() {
	var (
		refundPaymentInfo = model.RefundPaymentInfo{
			PaymentID: gofakeit.UUID(),
			Reason:    gofakeit.Word(),
		}
	)

	s.paymentRepository.On("Refund", s.ctx, refundPaymentInfo).Return(model.Payment{}, nil)

	res, err := s.service.RefundPayment(s.ctx, refundPaymentInfo)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *ServiceSuite) TestRefundPaymentRepoError() {
	var (
		repoErr = gofakeit.Error()

		refundPaymentInfo = model.RefundPaymentInfo{
			PaymentID: gofakeit.UUID(),
			Reason:    gofakeit.Word(),
		}
	)

	s.paymentRepository.On("Refund", s.ctx, refundPaymentInfo).Return(model.Payment{}, repoErr)

	res, err := s.service.RefundPayment(s.ctx, refundPaymentInfo)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(res)
}
