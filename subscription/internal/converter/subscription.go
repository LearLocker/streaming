package converter

import (
	subscriptionV1 "github.com/LearLocker/streaming/shared/pkg/openapi/subscription/v1"
	"github.com/LearLocker/streaming/subscription/internal/model"
	"github.com/google/uuid"
)

func CreateSubscriptionInfoFromAPI(req *subscriptionV1.CreateSubscriptionRequest) model.CreateSubscriptionInfo {
	return model.CreateSubscriptionInfo{
		PlanID:        req.PlanID.String(),
		PaymentMethod: model.PaymentMethod(req.PaymentMethod),
	}
}

func PaySubscriptionInfoFromAPI(
	uuid string,
	req *subscriptionV1.PaySubscriptionRequest,
) model.PaySubscriptionInfo {
	return model.PaySubscriptionInfo{
		UUID:          uuid,
		PaymentMethod: model.PaymentMethod(req.PaymentMethod),
	}
}

func SubscriptionToAPI(sub *model.Subscription) *subscriptionV1.Subscription {
	resp := &subscriptionV1.Subscription{
		UUID:          subscriptionV1.NewOptUUID(mustParseUUID(sub.UUID)),
		Status:        subscriptionV1.NewOptSubscriptionStatus(subscriptionV1.SubscriptionStatus(sub.Status)),
		PlanID:        subscriptionV1.NewOptString(sub.PlanID),
		PlanName:      subscriptionV1.NewOptString(sub.PlanName),
		Amount:        subscriptionV1.NewOptInt64(sub.Amount),
		Currency:      subscriptionV1.NewOptString(sub.Currency),
		PaymentMethod: subscriptionV1.NewOptPaymentMethod(subscriptionV1.PaymentMethod(sub.PaymentMethod)),
	}

	if sub.ExpiresAt != nil {
		resp.ExpiresAt = subscriptionV1.NewOptDateTime(*sub.ExpiresAt)
	}

	return resp
}

func PaySubscriptionResponseToAPI(sub *model.Subscription) *subscriptionV1.PaySubscriptionResponse {
	resp := &subscriptionV1.PaySubscriptionResponse{
		UUID:          subscriptionV1.NewOptUUID(mustParseUUID(sub.UUID)),
		Status:        subscriptionV1.NewOptSubscriptionStatus(subscriptionV1.SubscriptionStatus(sub.Status)),
		PaymentMethod: subscriptionV1.NewOptPaymentMethod(subscriptionV1.PaymentMethod(sub.PaymentMethod)),
	}

	if sub.ExpiresAt != nil {
		resp.ExpiresAt = subscriptionV1.NewOptDateTime(*sub.ExpiresAt)
	}

	return resp
}

func mustParseUUID(s string) uuid.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}
	return id
}
