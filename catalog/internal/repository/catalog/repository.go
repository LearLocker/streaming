package catalog

import (
	"context"
	"errors"
	"fmt"

	"github.com/LearLocker/streaming/catalog/internal/model"
	"github.com/LearLocker/streaming/catalog/internal/repository/converter"
	"github.com/LearLocker/streaming/catalog/internal/storage/mongoDb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	def "github.com/LearLocker/streaming/catalog/internal/repository"
)

var _ def.CatalogRepository = (*repository)(nil)

type repository struct {
	mongoClient *mongo.Client
}

func NewRepository(client *mongo.Client) *repository {
	return &repository{mongoClient: client}
}

func (r *repository) GetMovies(ctx context.Context, genreNames []string, page int32, pageSize int32) ([]*model.Movie, int32, error) {
	collection := mongoDb.OpenCollection(r.mongoClient, "movies")

	findOpts := options.Find()

	// пагинация
	if pageSize > 0 {
		findOpts.SetLimit(int64(pageSize))
	}
	if page > 1 && pageSize > 0 {
		findOpts.SetSkip(int64((page - 1) * pageSize))
	}

	// фильтр по жанрам если передан
	filter := bson.D{}
	if len(genreNames) > 0 {
		filter = bson.D{{
			Key:   "genre.genre_name",
			Value: bson.D{{Key: "$in", Value: genreNames}},
		}}
	}

	cursor, err := collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, 0, fmt.Errorf("find movies: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, 0, fmt.Errorf("decode movies: %w", err)
	}

	movies := make([]*model.Movie, 0, len(docs))
	for _, doc := range docs {
		movies = append(movies, converter.MovieFromBson(doc))
	}

	return movies, int32(len(movies)), nil
}

func (r *repository) GetMovieById(ctx context.Context, imdbId string) (*model.Movie, error) {
	collection := mongoDb.OpenCollection(r.mongoClient, "movies")

	var doc bson.M
	err := collection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: imdbId}}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("movie %s not found", imdbId)
		}
		return nil, fmt.Errorf("find movie: %w", err)
	}

	return converter.MovieFromBson(doc), nil
}

func (r *repository) AddMovie(ctx context.Context, info model.AddMovie) (string, error) {
	doc := bson.M{
		"imdb_id":     info.ImdbID,
		"title":       info.Title,
		"poster_path": info.PosterPath,
		"youtube_id":  info.YouTubeID,
		"genre":       converter.GenresToBson(info.Genre),
		"ranking": bson.M{
			"ranking_value": info.Ranking.RankingValue,
			"ranking_name":  info.Ranking.RankingName,
		},
	}

	collection := mongoDb.OpenCollection(r.mongoClient, "movies")
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return "", fmt.Errorf("insert movie: %w", err)
	}

	return result.InsertedID.(bson.ObjectID).Hex(), nil
}

func (r *repository) UpdateReview(ctx context.Context, imdbId string, authorId string, text string, ranking model.Ranking) error {
	filter := bson.D{{Key: "imdb_id", Value: imdbId}}
	update := bson.M{
		"$set": bson.M{
			"review": bson.M{
				"text":      text,
				"author_id": authorId,
			},
			"ranking": bson.M{
				"ranking_name":  ranking.RankingName,
				"ranking_value": ranking.RankingValue,
			},
		},
	}

	collection := mongoDb.OpenCollection(r.mongoClient, "movies")
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("update review: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("movie %s not found", imdbId)
	}

	return nil
}

func (r *repository) GetRecommendedMovies(ctx context.Context, userId string, limit int32) ([]*model.Movie, error) {
	genres, err := r.GetUserFavouriteGenres(ctx, userId)
	if err != nil {
		return nil, err
	}

	findOpts := options.Find()
	findOpts.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})
	if limit > 0 {
		findOpts.SetLimit(int64(limit))
	}

	filter := bson.D{{
		Key:   "genre.genre_name",
		Value: bson.D{{Key: "$in", Value: genres}},
	}}

	collection := mongoDb.OpenCollection(r.mongoClient, "movies")
	cursor, err := collection.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, fmt.Errorf("find recommended movies: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode recommended movies: %w", err)
	}

	movies := make([]*model.Movie, 0, len(docs))
	for _, doc := range docs {
		movies = append(movies, converter.MovieFromBson(doc))
	}

	return movies, nil
}

func (r *repository) GetGenres(ctx context.Context) ([]*model.Genre, error) {
	collection := mongoDb.OpenCollection(r.mongoClient, "genres")

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("find genres: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode genres: %w", err)
	}

	genres := make([]*model.Genre, 0, len(docs))
	for _, doc := range docs {
		genres = append(genres, &model.Genre{
			GenreID:   converter.ToInt(doc["genre_id"]),
			GenreName: converter.ToString(doc["genre_name"]),
		})
	}

	return genres, nil
}

func (r *repository) GetPlan(ctx context.Context, planId string) (*model.Plan, error) {
	collection := mongoDb.OpenCollection(r.mongoClient, "plans")

	var doc bson.M
	err := collection.FindOne(ctx, bson.M{"plan_id": planId}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("plan %s not found", planId)
		}
		return nil, fmt.Errorf("find plan: %w", err)
	}

	return converter.PlanFromBson(doc), nil
}

func (r *repository) ListPlans(ctx context.Context, onlyActive bool) ([]*model.Plan, error) {
	filter := bson.M{}
	if onlyActive {
		filter["is_active"] = true
	}

	collection := mongoDb.OpenCollection(r.mongoClient, "plans")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("find plans: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode plans: %w", err)
	}

	plans := make([]*model.Plan, 0, len(docs))
	for _, doc := range docs {
		plans = append(plans, converter.PlanFromBson(doc))
	}

	return plans, nil
}

func (r *repository) GetUserFavouriteGenres(
	ctx context.Context,
	userID string,
) ([]string, error) {

	filter := bson.M{"user_id": userID}
	projection := bson.M{"favourite_genres.genre_name": 1, "_id": 0}

	collection := mongoDb.OpenCollection(r.mongoClient, "users")
	var result bson.M
	err := collection.FindOne(ctx, filter,
		options.FindOne().SetProjection(projection),
	).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("find user genres: %w", err)
	}

	favGenresArr, ok := result["favourite_genres"].(bson.A)
	if !ok {
		return []string{}, nil
	}

	var genreNames []string
	for _, g := range favGenresArr {
		if gm, ok := g.(bson.D); ok {
			for _, item := range gm {
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

func (r *repository) GetRankings(ctx context.Context) ([]*model.Ranking, error) {
	collection := mongoDb.OpenCollection(r.mongoClient, "ranking")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []bson.M
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("decode rankings: %w", err)
	}

	rankings := make([]*model.Ranking, 0, len(docs))
	for _, doc := range docs {
		rankings = append(rankings, &model.Ranking{
			RankingValue: converter.ToInt(doc["ranking_value"]),
			RankingName:  converter.ToString(doc["ranking_name"]),
		})
	}

	return rankings, nil
}
