package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello World!!",
		})
	})
	err := r.Run(":8080")
	if err != nil {
		return
	}
	fmt.Println("listening on port 8080")
}
