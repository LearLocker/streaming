package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/LearLocker/streaming/payment/internal/converter"
	"github.com/LearLocker/streaming/payment/internal/service"
	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}

func (a *api) ProcessPayment(
	ctx context.Context,
	req *paymentV1.ProcessPaymentRequest,
) (*paymentV1.ProcessPaymentResponse, error) {

	payment, err := a.paymentService.ProcessPayment(
		ctx,
		converter.ProcessPaymentInfoFromProto(req),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "process payment: %v", err)
	}

	return &paymentV1.ProcessPaymentResponse{
		PaymentId: payment.PaymentID,
		Status:    converter.PaymentStatusToProto(payment.Status),
		Message:   "payment processed successfully",
	}, nil
}

func (a *api) GetPayment(
	ctx context.Context,
	req *paymentV1.GetPaymentRequest,
) (*paymentV1.GetPaymentResponse, error) {

	payment, err := a.paymentService.GetPayment(ctx, req.GetPaymentId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get payment: %v", err)
	}

	return &paymentV1.GetPaymentResponse{
		Payment: converter.PaymentToProto(payment),
	}, nil
}

func (a *api) RefundPayment(
	ctx context.Context,
	req *paymentV1.RefundPaymentRequest,
) (*paymentV1.RefundPaymentResponse, error) {

	payment, err := a.paymentService.RefundPayment(
		ctx,
		converter.RefundPaymentInfoFromProto(req),
	)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "refund payment: %v", err)
	}

	return &paymentV1.RefundPaymentResponse{
		PaymentId:  payment.PaymentID,
		Status:     converter.PaymentStatusToProto(payment.Status),
		RefundedAt: converter.PaymentToProto(payment).RefundedAt,
	}, nil
}
