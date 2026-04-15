package routes

import (
	movie_controller "streaming/internal/http-server/handlers/movies"
	user_controller "streaming/internal/http-server/handlers/users"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupUnProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	router.GET("/movies", movie_controller.GetMovies(client))
	router.GET("/genres", movie_controller.GetGenres(client))

	router.POST("/user/register", user_controller.RegisterUser(client))
	router.POST("/user/login", user_controller.Login(client))
}
