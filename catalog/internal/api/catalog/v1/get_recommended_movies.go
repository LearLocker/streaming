package v1

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/converter"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) GetRecommendedMovies(ctx context.Context, req *catalogV1.GetRecommendedMoviesRequest) (*catalogV1.GetRecommendedMoviesResponse, error) {

	movies, err := a.catalogService.GetRecommendedMovies(
		ctx,
		req.GetUserId(),
		req.GetLimit(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get recommended movies: %v", err)
	}

	return &catalogV1.GetRecommendedMoviesResponse{
		Movies: converter.MoviesToProto(movies),
	}, nil
}
