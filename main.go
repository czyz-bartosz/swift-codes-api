package main

import (
	"awesomeProject/configs"
	"awesomeProject/dbs"
	"awesomeProject/migrations"
	"awesomeProject/routes"
	"github.com/uptrace/bun"
)

func main() {
	config := configs.GetConfig()

	db := dbs.Connect(
		config.DBConfig,
	)

	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)

	err := migrations.Migrate(db)
	if err != nil {
		panic(err)
	}

	router := routes.SetupRouter()
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
