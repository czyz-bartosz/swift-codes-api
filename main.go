package main

import (
	"awesomeProject/configs"
	"awesomeProject/dbs"
	"awesomeProject/internal/dbimporter/utils"
	"awesomeProject/migrations"
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

	bankRepo := &repositories.BankRepoPostgres{
		Db: db,
	}

	router := routes.SetupRouter(bankRepo)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
