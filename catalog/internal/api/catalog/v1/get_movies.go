package v1

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/converter"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) GetMovies(ctx context.Context, req *catalogV1.GetMoviesRequest) (*catalogV1.GetMoviesResponse, error) {
	movies, total, err := a.catalogService.GetMovies(
		ctx,
		converter.GetMoviesFilterFromProto(req),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get movies: %v", err)
	}

	return &catalogV1.GetMoviesResponse{
		Movies: converter.MoviesToProto(movies),
		Total:  total,
	}, nil
}
