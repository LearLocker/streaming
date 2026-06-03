package v1

import (
	"github.com/LearLocker/streaming/payment/internal/converter"
	"github.com/LearLocker/streaming/payment/internal/model"
	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *APISuite) TestRefundPaymentSuccess() {
	var (
		paymentID     = gofakeit.UUID()
		paymentStatus = converter.PaymentStatusToProto(model.StatusRefunded)
		refundAt      = timestamppb.New(gofakeit.Date())

		refundPaymentRequest = &paymentV1.RefundPaymentRequest{
			PaymentId: paymentID,
			Reason:    gofakeit.Word(),
		}
	)

	s.paymentService.On("RefundPayment", s.ctx, refundPaymentRequest).Return(paymentID, paymentStatus, refundAt, nil)

	res, err := s.api.RefundPayment(s.ctx, refundPaymentRequest)
	s.Require().NoError(err)
	s.Require().NotNil(res)

	s.Require().Equal(paymentID, res.GetPaymentId())
	s.Require().Equal(paymentStatus, res.GetStatus())
	s.Require().Equal(refundAt, res.GetRefundedAt())
}

func (s *APISuite) TestRefundPaymentNotFound() {
	var (
		refundPaymentRequest = &paymentV1.RefundPaymentRequest{
			PaymentId: gofakeit.UUID(),
			Reason:    gofakeit.Word(),
		}
	)

	s.paymentService.On("RefundPayment", s.ctx, refundPaymentRequest).Return(nil, model.ErrPaymentNotFound)

	res, err := s.api.RefundPayment(s.ctx, refundPaymentRequest)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())

}

func (s *APISuite) TestRefundPaymentServiceError() {
	var (
		serviceErr = gofakeit.Error()

		refundPaymentRequest = &paymentV1.RefundPaymentRequest{
			PaymentId: gofakeit.UUID(),
			Reason:    gofakeit.Word(),
		}
	)

	s.paymentService.On("RefundPayment", s.ctx, refundPaymentRequest).Return(nil, serviceErr)

	res, err := s.api.RefundPayment(s.ctx, refundPaymentRequest)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
