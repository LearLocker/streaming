package movies

import (
	"context"
	"net/http"
	"time"

	"streaming/internal/domain/models"
	"streaming/internal/storage/mongoDb"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
