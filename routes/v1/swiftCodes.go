package v1

import (
	"awesomeProject/controllers"
	"awesomeProject/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupGroup(group *gin.RouterGroup, bankRepo repositories.BankRepo) {
	controller := controllers.Controller{
		BankRepo: bankRepo,
	}

	group.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World v1",
		})
	})

	group.GET("/:swiftCode", controller.GetDetails)
}
