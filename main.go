package main

import (
	"awesomeProject/configs"
	"awesomeProject/dbs"
	"awesomeProject/routes"
)

func main() {
	config := configs.GetConfig()

	_ = dbs.Connect(
		config.DBConfig,
	)

	router := routes.SetupRouter()
	err := router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
