package v1

import (
	"context"
	"errors"
	"net/http"

	subscriptionV1 "github.com/LearLocker/streaming/shared/pkg/openapi/subscription/v1"
	"github.com/LearLocker/streaming/subscription/internal/clients"
	"github.com/LearLocker/streaming/subscription/internal/converter"
	"github.com/LearLocker/streaming/subscription/internal/service"
)

type api struct {
	subscriptionV1.UnimplementedHandler

	subscriptionService service.SubscriptionService
}

func NewAPI(subscriptionService service.SubscriptionService) *api {
	return &api{subscriptionService: subscriptionService}
}

func (a *api) NewError(
	_ context.Context,
	err error,
) *subscriptionV1.GenericErrorStatusCode {
	return &subscriptionV1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: subscriptionV1.GenericError{
			Code:    subscriptionV1.NewOptInt(http.StatusInternalServerError),
			Message: subscriptionV1.NewOptString(err.Error()),
		},
	}
}

func (a *api) CreateSubscription(
	ctx context.Context,
	req *subscriptionV1.CreateSubscriptionRequest,
) (subscriptionV1.CreateSubscriptionRes, error) {

	sub, err := a.subscriptionService.CreateSubscription(
		ctx,
		converter.CreateSubscriptionInfoFromAPI(req),
	)
	if err != nil {
		if errors.Is(err, clients.ErrPlanNotFound) {
			return &subscriptionV1.NotFoundError{Message: err.Error()}, nil
		}
		return &subscriptionV1.InternalServerError{Message: err.Error()}, nil
	}

	return converter.SubscriptionToAPI(sub), nil
}

func (a *api) GetSubscriptionByUuid(
	ctx context.Context,
	params subscriptionV1.GetSubscriptionByUuidParams,
) (subscriptionV1.GetSubscriptionByUuidRes, error) {

	sub, err := a.subscriptionService.GetSubscription(ctx, params.UUID)
	if err != nil {
		return &subscriptionV1.NotFoundError{Message: err.Error()}, nil
	}

	return converter.SubscriptionToAPI(sub), nil
}

func (a *api) PaySubscriptionByUuid(
	ctx context.Context,
	req *subscriptionV1.PaySubscriptionRequest,
	params subscriptionV1.PaySubscriptionByUuidParams,
) (subscriptionV1.PaySubscriptionByUuidRes, error) {

	sub, err := a.subscriptionService.PaySubscription(
		ctx,
		converter.PaySubscriptionInfoFromAPI(params.UUID, req),
	)
	if err != nil {
		switch {
		case errors.Is(err, errAlreadyPaid), errors.Is(err, errCancelled):
			return &subscriptionV1.ConflictError{Message: err.Error()}, nil
		default:
			return &subscriptionV1.InternalServerError{Message: err.Error()}, nil
		}
	}

	return converter.PaySubscriptionResponseToAPI(sub), nil
}

func (a *api) CancelSubscriptionByUuid(
	ctx context.Context,
	params subscriptionV1.CancelSubscriptionByUuidParams,
) (subscriptionV1.CancelSubscriptionByUuidRes, error) {

	err := a.subscriptionService.CancelSubscription(ctx, params.UUID)
	if err != nil {
		if errors.Is(err, errAlreadyCancelled) {
			return &subscriptionV1.ConflictError{Message: err.Error()}, nil
		}
		return &subscriptionV1.NotFoundError{Message: err.Error()}, nil
	}

	return &subscriptionV1.CancelSubscriptionByUuidNoContent{}, nil
}

// sentinel ошибки для маппинга в HTTP коды
var (
	errAlreadyPaid      = errors.New("subscription is already paid")
	errCancelled        = errors.New("subscription is cancelled")
	errAlreadyCancelled = errors.New("subscription is already cancelled")
)
