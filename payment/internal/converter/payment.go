package converter

import (
	"github.com/LearLocker/streaming/payment/internal/model"
	paymentV1 "github.com/LearLocker/streaming/shared/pkg/proto/payment/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProcessPaymentInfoFromProto(req *paymentV1.ProcessPaymentRequest) model.ProcessPaymentInfo {
	return model.ProcessPaymentInfo{
		SubscriptionID: req.GetSubscriptionId(),
		UserID:         req.GetUserId(),
		Amount:         req.GetAmount(),
		Currency:       req.GetCurrency(),
		Method:         req.GetMethod(),
	}
}

func RefundPaymentInfoFromProto(req *paymentV1.RefundPaymentRequest) model.RefundPaymentInfo {
	return model.RefundPaymentInfo{
		PaymentID: req.GetPaymentId(),
		Reason:    req.GetReason(),
	}
}

func PaymentToProto(p *model.Payment) *paymentV1.Payment {
	proto := &paymentV1.Payment{
		PaymentId:      p.PaymentID,
		SubscriptionId: p.SubscriptionID,
		UserId:         p.UserID,
		Amount:         p.Amount,
		Currency:       p.Currency,
		Method:         p.Method,
		Status:         PaymentStatusToProto(p.Status),
		Reason:         p.Reason,
		CreatedAt:      timestamppb.New(p.CreatedAt),
	}

	if p.RefundedAt != nil {
		proto.RefundedAt = timestamppb.New(*p.RefundedAt)
	}

	return proto
}

func PaymentStatusToProto(s model.PaymentStatus) paymentV1.PaymentStatus {
	switch s {
	case model.StatusPending:
		return paymentV1.PaymentStatus_PAYMENT_STATUS_PENDING
	case model.StatusProcessing:
		return paymentV1.PaymentStatus_PAYMENT_STATUS_PROCESSING
	case model.StatusSuccess:
		return paymentV1.PaymentStatus_PAYMENT_STATUS_SUCCESS
	case model.StatusFailed:
		return paymentV1.PaymentStatus_PAYMENT_STATUS_FAILED
	case model.StatusRefunded:
		return paymentV1.PaymentStatus_PAYMENT_STATUS_REFUNDED
	default:
		return paymentV1.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}
