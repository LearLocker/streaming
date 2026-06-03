package v1

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/converter"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) GetGenres(ctx context.Context, req *catalogV1.GetGenresRequest) (*catalogV1.GetGenresResponse, error) {

	genres, err := a.catalogService.GetGenres(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get genres: %v", err)
	}

	return &catalogV1.GetGenresResponse{
		Genres: converter.GenresToProto(genres),
	}, nil
}
