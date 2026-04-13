package routes

import (
	movie_controller "streaming/internal/http-server/handlers/movies"
	"streaming/internal/http-server/middleware/auth"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.Use(auth.AuthMiddleWare())

	router.GET("/movie/:imdb_id", movie_controller.GetMovieById(client))
	router.POST("/movie", movie_controller.AddMovie(client))

}
