package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/LearLocker/streaming/internal/storage/mongoDb"
	catalogV1 "github.com/LearLocker/streaming/shared/pkg/proto/catalog/v1"
)

const grpcPort = 50051

type catalogService struct {
	catalogV1.UnimplementedCatalogServiceServer

	mongoClient *mongo.Client
}

func newCatalogService(client *mongo.Client) *catalogService {
	return &catalogService{mongoClient: client}
}

func (s *catalogService) GetMovies(
	ctx context.Context,
	req *catalogV1.GetMoviesRequest,
) (*catalogV1.GetMoviesResponse, error) {
	collection := mongoDb.OpenCollection(s.mongoClient, "movies")

	findOpts := options.Find()

	// пагинация
	if req.PageSize > 0 {
		findOpts.SetLimit(int64(req.PageSize))
	}
	if req.Page > 1 && req.PageSize > 0 {
		findOpts.SetSkip(int64((req.Page - 1) * req.PageSize))
	}

	// фильтр по жанрам если передан
	filter := bson.D{}
	if len(req.GenreNames) > 0 {
		filter = bson.D{{
			Key:   "genre.genre_name",
			Value: bson.D{{Key: "$in", Value: req.GenreNames}},
		}}
	}

	cursor, err := collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch movies: %v", err)
	}
	defer cursor.Close(ctx)

	var movies []bson.M
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decode movies: %v", err)
	}

	return &catalogV1.GetMoviesResponse{
		Movies: toProtoMovies(movies),
		Total:  int32(len(movies)),
	}, nil
}
func (s *catalogService) GetMovieById(ctx context.Context, req *catalogV1.GetMovieByIdRequest) (*catalogV1.GetMovieByIdResponse, error) {
	if req.ImdbId == "" {
		return nil, status.Error(codes.InvalidArgument, "imdb_id is required")
	}

	collection := mongoDb.OpenCollection(s.mongoClient, "movies")

	var movie bson.M
	err := collection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: req.ImdbId}}).Decode(&movie)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "movie %s not found", req.ImdbId)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch movie: %v", err)
	}

	return &catalogV1.GetMovieByIdResponse{
		Movie: toProtoMovie(movie),
	}, nil
}
func (s *catalogService) AddMovie(ctx context.Context, req *catalogV1.AddMovieRequest) (*catalogV1.AddMovieResponse, error) {
	if req.ImdbId == "" || req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "imdb_id and title are required")
	}

	doc := bson.M{
		"imdb_id":     req.ImdbId,
		"title":       req.Title,
		"poster_path": req.PosterPath,
		"youtube_id":  req.YoutubeId,
		"genre":       toProtoGenresBson(req.Genre),
		"ranking": bson.M{
			"ranking_value": req.Ranking.GetRankingValue(),
			"ranking_name":  req.Ranking.GetRankingName(),
		},
	}

	collection := mongoDb.OpenCollection(s.mongoClient, "movies")
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add movie: %v", err)
	}

	return &catalogV1.AddMovieResponse{
		InsertedId: result.InsertedID.(bson.ObjectID).Hex(),
	}, nil
}
func (s *catalogService) UpdateReview(ctx context.Context, req *catalogV1.UpdateReviewRequest) (*catalogV1.UpdateReviewResponse, error) {
	if req.ImdbId == "" {
		return nil, status.Error(codes.InvalidArgument, "imdb_id is required")
	}
	if req.Text == "" {
		return nil, status.Error(codes.InvalidArgument, "review text is required")
	}

	// получить sentiment через LLM — логика из GetReviewRanking()
	sentiment, rankVal, err := s.getReviewRanking(ctx, req.Text)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get review ranking: %v", err)
	}

	filter := bson.D{{Key: "imdb_id", Value: req.ImdbId}}
	update := bson.M{
		"$set": bson.M{
			"review": bson.M{
				"text":      req.Text,
				"author_id": req.AuthorId,
			},
			"ranking": bson.M{
				"ranking_name":  sentiment,
				"ranking_value": rankVal,
			},
		},
	}

	collection := mongoDb.OpenCollection(s.mongoClient, "movies")
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update review: %v", err)
	}
	if result.MatchedCount == 0 {
		return nil, status.Errorf(codes.NotFound, "movie %s not found", req.ImdbId)
	}

	return &catalogV1.UpdateReviewResponse{
		RankingName:  sentiment,
		RankingValue: int32(rankVal),
		Text:         req.Text,
	}, nil
}

func (s *catalogService) getReviewRanking(ctx context.Context, review string) (string, int, error) {
	rankings, err := s.getRankings(ctx)
	if err != nil {
		return "", 0, err
	}

	var parts []string
	for _, r := range rankings {
		if r["ranking_value"] != 999 {
			if name, ok := r["ranking_name"].(string); ok {
				parts = append(parts, name)
			}
		}
	}
	sentimentDelimited := strings.Join(parts, ",")

	_ = godotenv.Load("internal/config/.env")

	openAiApiKey := os.Getenv("OPEN_AI_API_KEY")
	if openAiApiKey == "" {
		return "", 0, errors.New("OPEN_AI_API_KEY is not set")
	}

	basePromptTemp := os.Getenv("BASE_PROMPT_TEMP")
	if basePromptTemp == "" {
		return "", 0, errors.New("BASE_PROMPT_TEMP is not set")
	}

	llm, err := openai.New(openai.WithToken(openAiApiKey))
	if err != nil {
		return "", 0, err
	}

	basePrompt := strings.Replace(basePromptTemp, "{rankings}", sentimentDelimited, 1)
	response, err := llm.Call(ctx, basePrompt+review)
	if err != nil {
		return "", 0, err
	}

	// найти ranking_value по имени из ответа LLM
	rankVal := 0
	for _, r := range rankings {
		if name, ok := r["ranking_name"].(string); ok && name == response {
			if val, ok := r["ranking_value"].(int32); ok {
				rankVal = int(val)
			}
			break
		}
	}

	return response, rankVal, nil
}

func (s *catalogService) getRankings(ctx context.Context) ([]bson.M, error) {
	collection := mongoDb.OpenCollection(s.mongoClient, "ranking")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var rankings []bson.M
	if err = cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

func (s *catalogService) GetRecommendedMovies(ctx context.Context, req *catalogV1.GetRecommendedMoviesRequest) (*catalogV1.GetRecommendedMoviesResponse, error) {
	favoriteGenres, err := s.getUserFavouriteGenres(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get favourite genres: %v", err)
	}

	limit := req.Limit
	if limit == 0 {
		// fallback на env переменную как в оригинале
		_ = godotenv.Load("internal/config/.env")
		if v := os.Getenv("RECOMMENDED_MOVIE_LIMIT"); v != "" {
			if n, err := strconv.ParseInt(v, 10, 32); err == nil {
				limit = int32(n)
			}
		}
	}

	findOpts := options.Find()
	findOpts.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
	if limit > 0 {
		findOpts.SetLimit(int64(limit))
	}

	filter := bson.D{{
		Key:   "genre.genre_name",
		Value: bson.D{{Key: "$in", Value: favoriteGenres}},
	}}

	collection := mongoDb.OpenCollection(s.mongoClient, "movies")
	cursor, err := collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch recommended movies: %v", err)
	}
	defer cursor.Close(ctx)

	var movies []bson.M
	if err = cursor.All(ctx, &movies); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decode movies: %v", err)
	}

	return &catalogV1.GetRecommendedMoviesResponse{
		Movies: toProtoMovies(movies),
	}, nil
}

func (s *catalogService) getUserFavouriteGenres(ctx context.Context, userID string) ([]string, error) {
	filter := bson.M{"user_id": userID}
	projection := bson.M{
		"favourite_genres.genre_name": 1,
		"_id":                         0,
	}

	opts := options.FindOne().SetProjection(projection)
	var result bson.M

	collection := mongoDb.OpenCollection(s.mongoClient, "users")
	err := collection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []string{}, nil
		}
		return nil, err
	}

	favGenresArr, ok := result["favourite_genres"].(bson.A)
	if !ok {
		return []string{}, nil
	}

	var genreNames []string
	for _, favGenre := range favGenresArr {
		if genreMap, ok := favGenre.(bson.D); ok {
			for _, item := range genreMap {
				if item.Key == "genre_name" {
					if name, ok := item.Value.(string); ok {
						genreNames = append(genreNames, name)
					}
				}
			}
		}
	}

	return genreNames, nil
}

func (s *catalogService) GetGenres(ctx context.Context, req *catalogV1.GetGenresRequest) (*catalogV1.GetGenresResponse, error) {
	collection := mongoDb.OpenCollection(s.mongoClient, "genres")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch genres: %v", err)
	}
	defer cursor.Close(ctx)

	var genres []bson.M
	if err = cursor.All(ctx, &genres); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decode genres: %v", err)
	}

	var protoGenres []*catalogV1.Genre
	for _, g := range genres {
		protoGenres = append(protoGenres, &catalogV1.Genre{
			GenreId:   int32(toInt(g["genre_id"])),
			GenreName: toString(g["genre_name"]),
		})
	}

	return &catalogV1.GetGenresResponse{Genres: protoGenres}, nil
}
func (s *catalogService) GetPlan(ctx context.Context, req *catalogV1.GetPlanRequest) (*catalogV1.GetPlanResponse, error) {
	if req.PlanId == "" {
		return nil, status.Error(codes.InvalidArgument, "plan_id is required")
	}

	collection := mongoDb.OpenCollection(s.mongoClient, "plans")

	var plan bson.M
	err := collection.FindOne(ctx, bson.M{"plan_id": req.PlanId}).Decode(&plan)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, status.Errorf(codes.NotFound, "plan %s not found", req.PlanId)
		}
		return nil, status.Errorf(codes.Internal, "failed to fetch plan: %v", err)
	}

	return &catalogV1.GetPlanResponse{
		Plan: toProtoPlan(plan),
	}, nil
}
func (s *catalogService) ListPlans(ctx context.Context, req *catalogV1.ListPlansRequest) (*catalogV1.ListPlansResponse, error) {
	filter := bson.M{}
	if req.OnlyActive {
		filter["is_active"] = true
	}

	collection := mongoDb.OpenCollection(s.mongoClient, "plans")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch plans: %v", err)
	}
	defer cursor.Close(ctx)

	var plans []bson.M
	if err = cursor.All(ctx, &plans); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to decode plans: %v", err)
	}

	var protoPlans []*catalogV1.Plan
	for _, p := range plans {
		protoPlans = append(protoPlans, toProtoPlan(p))
	}

	return &catalogV1.ListPlansResponse{Plans: protoPlans}, nil
}

func toProtoMovies(movies []bson.M) []*catalogV1.Movie {
	result := make([]*catalogV1.Movie, 0, len(movies))
	for _, m := range movies {
		result = append(result, toProtoMovie(m))
	}
	return result
}

func toProtoMovie(m bson.M) *catalogV1.Movie {
	movie := &catalogV1.Movie{
		ImdbId:     toString(m["imdb_id"]),
		Title:      toString(m["title"]),
		PosterPath: toString(m["poster_path"]),
		YoutubeId:  toString(m["youtube_id"]),
	}

	if genres, ok := m["genre"].(bson.A); ok {
		for _, g := range genres {
			if gm, ok := g.(bson.M); ok {
				movie.Genre = append(movie.Genre, &catalogV1.Genre{
					GenreId:   int32(toInt(gm["genre_id"])),
					GenreName: toString(gm["genre_name"]),
				})
			}
		}
	}

	if ranking, ok := m["ranking"].(bson.M); ok {
		movie.Ranking = &catalogV1.Ranking{
			RankingValue: int32(toInt(ranking["ranking_value"])),
			RankingName:  toString(ranking["ranking_name"]),
		}
	}

	return movie
}

func toProtoPlan(p bson.M) *catalogV1.Plan {
	return &catalogV1.Plan{
		PlanId:       toString(p["plan_id"]),
		Name:         toString(p["name"]),
		Price:        toInt64(p["price"]),
		Currency:     toString(p["currency"]),
		DurationDays: int32(toInt(p["duration_days"])),
		IsActive:     toBool(p["is_active"]),
	}
}

func toProtoGenresBson(genres []*catalogV1.Genre) bson.A {
	result := make(bson.A, 0, len(genres))
	for _, g := range genres {
		result = append(result, bson.M{
			"genre_id":   g.GenreId,
			"genre_name": g.GenreName,
		})
	}
	return result
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toInt(v any) int {
	switch n := v.(type) {
	case int:
		return n
	case int32:
		return int(n)
	case int64:
		return int(n)
	case float64:
		return int(n) // MongoDB числа часто приходят как float64
	}
	return 0
}

func toInt64(v any) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case int32:
		return int64(n)
	case float64:
		return int64(n)
	}
	return 0
}

func toBool(v any) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	mongoClient := mongoDb.Connect()

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем наш сервис
	catalogV1.RegisterCatalogServiceServer(grpcServer, newCatalogService(mongoClient))

	// Включаем рефлексию для отладки
	reflection.Register(grpcServer)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = grpcServer.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("✅ Server stopped")
}
