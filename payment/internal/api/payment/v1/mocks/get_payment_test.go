package v1

import (
	"github.com/LearLocker/streaming/payment/internal/converter"
	"github.com/LearLocker/streaming/payment/internal/model"
	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestGetPaymentSuccess() {
	var (
		paymentId      = gofakeit.UUID()
		subscriptionId = gofakeit.UUID()
		userID         = gofakeit.UUID()
		price          = gofakeit.Price(12, 100)
		currency       = gofakeit.CurrencyShort()

		getPaymentRequest = &paymentV1.GetPaymentRequest{
			PaymentId: paymentId,
		}

		payment = &model.Payment{
			PaymentID:      paymentId,
			SubscriptionID: subscriptionId,
			UserID:         userID,
			Amount:         int64(price),
			Currency:       currency,
			Method:         gofakeit.Word(),
			Status:         RandomPaymentStatus(),
			Reason:         gofakeit.Word(),
		}

		expectedProtoPayment = converter.PaymentToProto(payment)
	)
	s.paymentService.On("GetPayment", s.ctx, getPaymentRequest).Return(expectedProtoPayment, nil)

	res, err := s.api.GetPayment(s.ctx, getPaymentRequest)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedProtoPayment, res.GetPayment())

}

func (s *APISuite) TestGetPaymentNotFound() {
	var (
		getPaymentRequest = &paymentV1.GetPaymentRequest{
			PaymentId: gofakeit.UUID(),
		}
	)

	s.paymentService.On("GetPayment", s.ctx, getPaymentRequest).Return(nil, model.ErrPaymentNotFound)

	res, err := s.api.GetPayment(s.ctx, getPaymentRequest)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetPaymentServiceError() {
	var (
		serviceErr = gofakeit.Error()

		getPaymentRequest = &paymentV1.GetPaymentRequest{
			PaymentId: gofakeit.UUID(),
		}
	)

	s.paymentService.On("GetPayment", s.ctx, getPaymentRequest).Return(nil, serviceErr)

	res, err := s.api.GetPayment(s.ctx, getPaymentRequest)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
