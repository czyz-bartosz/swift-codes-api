package routes

import (
	"awesomeProject/repositories"
	"awesomeProject/routes/v1"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SetupRouter(swiftRepo repositories.SwiftRepo, validate *validator.Validate) *gin.Engine {
	router := gin.Default()

	v1.SetupGroup(router.Group("/v1/swift-codes"), swiftRepo, validate)

	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}

	return router
}
