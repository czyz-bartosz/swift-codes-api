package v1

import (
	"awesomeProject/controllers"
	"github.com/gin-gonic/gin"
)

func SetupGroup(group *gin.RouterGroup, controller *controllers.Controller) {

	group.POST("/", controller.AddSwift)
	group.GET("/:swiftCode", controller.GetSwiftDetails)
	group.DELETE("/:swiftCode", controller.DeleteSwift)
	group.GET("/country/:countryIso2Code", controller.GetSwiftsDetailsByCountryIso2Code)
}
