package v1

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/converter"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) GetMovieById(ctx context.Context, req *catalogV1.GetMovieByIdRequest) (*catalogV1.GetMovieByIdResponse, error) {
	movie, err := a.catalogService.GetMovieById(ctx, req.GetImdbId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get movie: %v", err)
	}

	return &catalogV1.GetMovieByIdResponse{
		Movie: converter.MovieToProto(movie),
	}, nil
}
