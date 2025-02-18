package main

import (
	"awesomeProject/configs"
	"awesomeProject/dbs"
	"awesomeProject/dbs/migrations"
	"awesomeProject/internal/dbimporter/utils"
	"awesomeProject/repositories"
	"awesomeProject/routes"
	"fmt"
	"github.com/uptrace/bun"
)

func main() {
	config := configs.GetConfig()

	db := dbs.Connect(
		&config.DBConfig,
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

	err = utils.ImportData("./data.csv")
	if err != nil {
		fmt.Println(err)
	}

	swiftRepo := &repositories.SwiftRepoPostgres{
		Db: db,
	}

	router := routes.SetupRouter(swiftRepo)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
