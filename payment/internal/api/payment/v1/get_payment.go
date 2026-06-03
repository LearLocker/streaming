package v1

import (
	"context"

	"github.com/LearLocker/streaming/payment/internal/converter"
	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) GetPayment(
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
