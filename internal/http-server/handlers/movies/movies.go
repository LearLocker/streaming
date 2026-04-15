package movies

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"streaming/internal/lib/jwt"
	"strings"
	"time"

	"streaming/internal/domain/models"
	"streaming/internal/storage/mongoDb"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func GetMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieCollection := mongoDb.OpenCollection(client, "movies")

		cursor, err := movieCollection.Find(ctx, bson.D{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		}
		defer cursor.Close(ctx)
		var movies []models.Movie

		if err := cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		}

		c.JSON(http.StatusOK, gin.H{"movies": movies})
	}
}

func GetMovieById(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
			return
		}

		movieCollection := mongoDb.OpenCollection(client, "movies")

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: movieID}}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"movie": movie})
	}
}

func AddMovie(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		localValidator := validator.New()
		if err := localValidator.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		}

		movieCollection := mongoDb.OpenCollection(client, "movies")
		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
		}

		c.JSON(http.StatusCreated, result)

	}
}

func ReviewUpdate(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie id is required"})
			return
		}

		userId, err := jwt.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		var req struct {
			Review string `json:"review"`
		}
		var resp struct {
			RankingName string `json:"ranking_name"`
			Review      string `json:"review"`
		}

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		sentiment, rankVal, err := GetReviewRanking(req.Review, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting review ranking"})
			return
		}

		filter := bson.D{{Key: "imdb_id", Value: movieId}}

		update := bson.M{
			"$set": bson.M{
				"review": bson.M{
					"text":      req.Review,
					"author_id": userId,
				},
				"ranking": bson.M{
					"ranking_name":  sentiment,
					"ranking_value": rankVal,
				},
			},
		}

		movieCollection := mongoDb.OpenCollection(client, "movies")

		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating movie review"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		resp.RankingName = sentiment
		resp.Review = req.Review

		c.JSON(http.StatusOK, resp)
	}
}

func GetReviewRanking(review string, client *mongo.Client, c *gin.Context) (string, int, error) {
	rankings, err := GetRankings(client)
	if err != nil {
		return "", 0, err
	}

	sentimentDelimited := ""

	for _, ranking := range rankings {
		if ranking.RankingValue != 999 {
			sentimentDelimited += ranking.RankingName + ","
		}
	}

	sentimentDelimited = strings.Trim(sentimentDelimited, ",")

	err = godotenv.Load("internal/config/.env")
	if err != nil {
		log.Println("Error loading .env file")
	}

	openAiApiKey := os.Getenv("OPEN_AI_API_KEY")
	if openAiApiKey == "" {
		return "", 0, errors.New("could not read OPEN_AI_API_KEY")
	}

	llm, err := openai.New(openai.WithToken(openAiApiKey))
	if err != nil {
		return "", 0, err
	}

	basePromptTemp := os.Getenv("BASE_PROMPT_TEMP")
	if basePromptTemp == "" {
		return "", 0, errors.New("prompt could not be empty")
	}

	basePrompt := strings.Replace(basePromptTemp, "{rankings}", sentimentDelimited, 1)

	response, err := llm.Call(c, basePrompt+review)
	if err != nil {
		return "", 0, err
	}

	rankVal := 0

	for _, ranking := range rankings {
		if ranking.RankingName == response {
			rankVal = ranking.RankingValue
			break
		}
	}

	return response, rankVal, nil
}

func GetRankings(client *mongo.Client) ([]models.Ranking, error) {
	var rankings []models.Ranking

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	rankingCollection := mongoDb.OpenCollection(client, "ranking")
	cursor, err := rankingCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}

	return rankings, nil
}

func GetRecommendedMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		userId, err := jwt.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		favoriteGenres, err := GetUsersFavouriteGenres(userId, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = godotenv.Load("internal/config/.env")
		if err != nil {
			log.Println("Error loading .env file")
		}

		var recommendedMovieLimitVal int64
		recommendedMovieLimitStr := os.Getenv("RECOMMENDED_MOVIE_LIMIT")

		if recommendedMovieLimitStr != "" {
			recommendedMovieLimitVal, _ = strconv.ParseInt(recommendedMovieLimitStr, 10, 64)
		}

		opts := options.Find()

		opts.SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}})

		opts.SetLimit(recommendedMovieLimitVal)

		filter := bson.D{
			{Key: "genre.genre_name", Value: bson.D{
				{Key: "$in", Value: favoriteGenres},
			}},
		}

		movieCollection := mongoDb.OpenCollection(client, "movies")

		var recommendedMovies []models.Movie

		cursor, err := movieCollection.Find(ctx, filter, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, recommendedMovies)
	}
}

func GetUsersFavouriteGenres(userId string, client *mongo.Client, c *gin.Context) ([]string, error) {
	var ctx, cancel = context.WithTimeout(c, 100*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userId}

	projection := bson.M{
		"favorite_genres.genre_name": 1,
		"_id":                        0,
	}

	opts := options.FindOne().SetProjection(projection)
	var result bson.M

	userCollection := mongoDb.OpenCollection(client, "users")
	err := userCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []string{}, nil
		}

		return []string{}, errors.New("")
	}

	favGenresArr, ok := result["favorite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("unable to retrieve favorite genres for user")
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

func GetGenres(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(c, 100*time.Second)
		defer cancel()

		genreCollection := mongoDb.OpenCollection(client, "genres")

		cursor, err := genreCollection.Find(ctx, bson.D{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching movie genres"})
			return
		}
		defer cursor.Close(ctx)

		var genres []models.Genre
		if err := cursor.All(ctx, &genres); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, genres)

	}
}
