package v1

import (
	"awesomeProject/controllers"
	"awesomeProject/repositories"
	"github.com/gin-gonic/gin"
)

func SetupGroup(group *gin.RouterGroup, swiftRepo repositories.SwiftRepo) {
	controller := controllers.Controller{
		SwiftRepo: swiftRepo,
	}

	group.POST("/", controller.AddSwift)
	group.GET("/:swiftCode", controller.GetSwiftDetails)
	group.DELETE("/:swiftCode", controller.DeleteSwift)
	group.GET("/country/:countryIso2Code", controller.GetSwiftsDetailsByCountryIso2Code)
}
