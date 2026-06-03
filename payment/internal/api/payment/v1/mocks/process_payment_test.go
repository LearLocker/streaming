package v1

import (
	"github.com/LearLocker/streaming/payment/internal/converter"
	"github.com/LearLocker/streaming/payment/internal/model"
	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestProcessPaymentSuccess() {
	var (
		paymentID     = gofakeit.UUID()
		paymentStatus = converter.PaymentStatusToProto(model.StatusSuccess)
		msg           = gofakeit.Paragraph(3, 5, 5, " ")

		processPaymentRequest = &paymentV1.ProcessPaymentRequest{
			SubscriptionId: gofakeit.UUID(),
			UserId:         gofakeit.UUID(),
			Amount:         int64(gofakeit.Price(12, 100)),
			Currency:       gofakeit.CurrencyShort(),
			Method:         gofakeit.Word(),
		}
	)
	s.paymentService.On("ProcessPayment", s.ctx, processPaymentRequest).Return(paymentID, paymentStatus, msg, nil)

	res, err := s.api.ProcessPayment(s.ctx, processPaymentRequest)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	s.Require().Equal(paymentID, res.GetPaymentId())
	s.Require().Equal(paymentStatus, res.GetStatus())
	s.Require().Equal(msg, res.GetMessage())
}

func (s *APISuite) TestProcessPaymentNotFound() {
	var (
		processPaymentRequest = &paymentV1.ProcessPaymentRequest{
			SubscriptionId: gofakeit.UUID(),
			UserId:         gofakeit.UUID(),
			Amount:         int64(gofakeit.Price(12, 100)),
			Currency:       gofakeit.CurrencyShort(),
			Method:         gofakeit.Word(),
		}
	)

	s.paymentService.On("ProcessPayment", s.ctx, processPaymentRequest).Return(nil, model.ErrPaymentNotFound)

	res, err := s.api.ProcessPayment(s.ctx, processPaymentRequest)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestProcessPaymentServiceError() {
	var (
		serviceErr = gofakeit.Error()

		processPaymentRequest = &paymentV1.ProcessPaymentRequest{
			SubscriptionId: gofakeit.UUID(),
			UserId:         gofakeit.UUID(),
			Amount:         int64(gofakeit.Price(12, 100)),
			Currency:       gofakeit.CurrencyShort(),
			Method:         gofakeit.Word(),
		}
	)

	s.paymentService.On("ProcessPayment", s.ctx, processPaymentRequest).Return(nil, serviceErr)

	res, err := s.api.ProcessPayment(s.ctx, processPaymentRequest)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
