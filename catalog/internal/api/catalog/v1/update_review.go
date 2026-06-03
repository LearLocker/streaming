package v1

import (
	"context"

	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) UpdateReview(ctx context.Context, req *catalogV1.UpdateReviewRequest) (*catalogV1.UpdateReviewResponse, error) {
	result, err := a.catalogService.UpdateReview(
		ctx,
		req.GetImdbId(),
		req.GetAuthorId(),
		req.GetText(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update review: %v", err)
	}

	return &catalogV1.UpdateReviewResponse{
		RankingName:  result.RankingName,
		RankingValue: result.RankingValue,
		Text:         result.Text,
	}, nil
}
