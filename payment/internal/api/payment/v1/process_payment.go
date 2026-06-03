package v1

import (
	"context"

	"github.com/LearLocker/streaming/payment/internal/converter"
	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) ProcessPayment(
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
