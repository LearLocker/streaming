package main

import (
	"fmt"
	"streaming/internal/storage/mongoDb"

	movie_controller "streaming/internal/http-server/handlers/movies"
	user_controller "streaming/internal/http-server/handlers/users"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "We are on go")
	})

	client := mongoDb.Connect()

	router.GET("/movies", movie_controller.GetMovies(client))
	router.GET("/movie/:imdb_id", movie_controller.GetMovieById(client))
	router.POST("/movie", movie_controller.AddMovie(client))
	router.POST("/user/register", user_controller.RegisterUser(client))
	router.POST("/user/login", user_controller.Login(client))

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Faild to start server", err)
	}
}
