package v1

import (
	"context"

	"github.com/LearLocker/streaming/payment/internal/converter"
	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) RefundPayment(
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
