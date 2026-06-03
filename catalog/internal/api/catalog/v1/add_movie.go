package v1

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/converter"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *Api) AddMovie(ctx context.Context, req *catalogV1.AddMovieRequest) (*catalogV1.AddMovieResponse, error) {
	insertedID, err := a.catalogService.AddMovie(
		ctx,
		converter.AddMovieInfoFromProto(req),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "add movie: %v", err)
	}

	return &catalogV1.AddMovieResponse{InsertedId: insertedID}, nil
}
