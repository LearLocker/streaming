package v1

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/converter"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) GetPlan(ctx context.Context, req *catalogV1.GetPlanRequest) (*catalogV1.GetPlanResponse, error) {

	plan, err := a.catalogService.GetPlan(ctx, req.GetPlanId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get plan: %v", err)
	}

	return &catalogV1.GetPlanResponse{
		Plan: converter.PlanToProto(plan),
	}, nil
}
