package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SetupGroup(group *gin.RouterGroup) {
	group.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World v1",
		})
	})
}
