package v1

import (
	"awesomeProject/controllers"
	"awesomeProject/repositories"
	"github.com/gin-gonic/gin"
)

func SetupGroup(group *gin.RouterGroup, bankRepo repositories.BankRepo) {
	controller := controllers.Controller{
		BankRepo: bankRepo,
	}

	group.GET("/:swiftCode", controller.GetBankDetails)
	group.GET("/country/:countryIso2Code", controller.GetBanksDetailsByCountryIso2Code)
	group.POST("/", controller.AddBank)
}
