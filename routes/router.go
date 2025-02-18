package routes

import (
	"awesomeProject/repositories"
	"awesomeProject/routes/v1"
	"github.com/gin-gonic/gin"
)

func SetupRouter(swiftRepo repositories.SwiftRepo) *gin.Engine {
	router := gin.Default()

	v1.SetupGroup(router.Group("/v1/swift-codes"), swiftRepo)

	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}

	return router
}
