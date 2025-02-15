package routes

import (
	v1 "awesomeProject/routes/v1"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1.SetupGroup(router.Group("/v1"))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!!",
		})
	})

	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}

	return router
}
