package v1

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/converter"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) ListPlans(ctx context.Context, req *catalogV1.ListPlansRequest) (*catalogV1.ListPlansResponse, error) {

	plans, err := a.catalogService.ListPlans(ctx, req.GetOnlyActive())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list plans: %v", err)
	}

	return &catalogV1.ListPlansResponse{
		Plans: converter.PlansToProto(plans),
	}, nil
}
