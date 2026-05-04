package v1

import (
	"context"

	"github.com/LearLocker/streaming/catalog/internal/converter"
	"github.com/LearLocker/streaming/catalog/internal/service"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type api struct {
	catalogV1.UnimplementedCatalogServiceServer

	catalogService service.CatalogService
}

func NewAPI(catalogService service.CatalogService) *api {
	return &api{
		catalogService: catalogService,
	}
}

func (a *api) GetMovies(ctx context.Context, req *catalogV1.GetMoviesRequest) (*catalogV1.GetMoviesResponse, error) {
	movies, total, err := a.catalogService.GetMovies(
		ctx,
		req.GetGenreNames(),
		req.GetPage(),
		req.GetPageSize(),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get movies: %v", err)
	}

	return &catalogV1.GetMoviesResponse{
		Movies: converter.MoviesToProto(movies),
		Total:  total,
	}, nil
}

func (a *api) GetMovieById(ctx context.Context, req *catalogV1.GetMovieByIdRequest) (*catalogV1.GetMovieByIdResponse, error) {
	movie, err := a.catalogService.GetMovieById(ctx, req.GetImdbId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get movie: %v", err)
	}

	return &catalogV1.GetMovieByIdResponse{
		Movie: converter.MovieToProto(movie),
	}, nil
}

func (a *api) AddMovie(ctx context.Context, req *catalogV1.AddMovieRequest) (*catalogV1.AddMovieResponse, error) {
	insertedID, err := a.catalogService.AddMovie(
		ctx,
		converter.AddMovieInfoFromProto(req),
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "add movie: %v", err)
	}

	return &catalogV1.AddMovieResponse{InsertedId: insertedID}, nil
}

func (a *api) UpdateReview(ctx context.Context, req *catalogV1.UpdateReviewRequest) (*catalogV1.UpdateReviewResponse, error) {
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

func (a *api) GetRecommendedMovies(ctx context.Context, req *catalogV1.GetRecommendedMoviesRequest) (*catalogV1.GetRecommendedMoviesResponse, error) {

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

func (a *api) GetGenres(ctx context.Context, req *catalogV1.GetGenresRequest) (*catalogV1.GetGenresResponse, error) {

	genres, err := a.catalogService.GetGenres(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get genres: %v", err)
	}

	return &catalogV1.GetGenresResponse{
		Genres: converter.GenresToProto(genres),
	}, nil
}
func (a *api) GetPlan(ctx context.Context, req *catalogV1.GetPlanRequest) (*catalogV1.GetPlanResponse, error) {

	plan, err := a.catalogService.GetPlan(ctx, req.GetPlanId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "get plan: %v", err)
	}

	return &catalogV1.GetPlanResponse{
		Plan: converter.PlanToProto(plan),
	}, nil
}
func (a *api) ListPlans(ctx context.Context, req *catalogV1.ListPlansRequest) (*catalogV1.ListPlansResponse, error) {

	plans, err := a.catalogService.ListPlans(ctx, req.GetOnlyActive())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list plans: %v", err)
	}

	return &catalogV1.ListPlansResponse{
		Plans: converter.PlansToProto(plans),
	}, nil
}
