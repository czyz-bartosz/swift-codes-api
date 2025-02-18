package v1

import (
	"awesomeProject/controllers"
	"awesomeProject/repositories"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func SetupGroup(group *gin.RouterGroup, swiftRepo repositories.SwiftRepo, validate *validator.Validate) {
	controller := controllers.Controller{
		SwiftRepo: swiftRepo,
		Validate:  validate,
	}

	group.POST("/", controller.AddSwift)
	group.GET("/:swiftCode", controller.GetSwiftDetails)
	group.DELETE("/:swiftCode", controller.DeleteSwift)
	group.GET("/country/:countryIso2Code", controller.GetSwiftsDetailsByCountryIso2Code)
}
