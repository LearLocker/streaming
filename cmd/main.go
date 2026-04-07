package main

import (
	"fmt"
	"streaming/internal/storage/mongoDb"

	controller "streaming/internal/http-server/handlers/movies"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "We are on go")
	})

	client := mongoDb.Connect()

	router.GET("/movies", controller.GetMovies(client))
	router.GET("/movie/:imdb_id", controller.GetMovieById(client))
	router.POST("/movie", controller.AddMovie(client))

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Faild to start server", err)
	}
}
