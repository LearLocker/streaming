package main

import (
	"fmt"

	"streaming/internal/http-server/routes"
	"streaming/internal/storage/mongoDb"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "We are on go")
	})

	client := mongoDb.Connect()

	routes.SetupUnProtectedRoutes(router, client)
	routes.SetupProtectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Faild to start server", err)
	}
}
