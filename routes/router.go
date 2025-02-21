package routes

import (
	"awesomeProject/controllers"
	"awesomeProject/routes/v1"
	"github.com/gin-gonic/gin"
)

func SetupRouter(swiftController *controllers.Controller) *gin.Engine {
	router := gin.Default()

	v1.SetupGroup(router.Group("/v1/swift-codes"), swiftController)

	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}

	return router
}
