package main

import (
	"awesomeProject/configs"
	"awesomeProject/dbs"
	"awesomeProject/dbs/migrations"
	"awesomeProject/internal/dbimporter/utils"
	"awesomeProject/models"
	"awesomeProject/repositories"
	"awesomeProject/routes"
	"fmt"
	"github.com/go-playground/validator/v10"
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

	validate := validator.New()
	validate.RegisterStructValidation(models.SwiftStructLevelValidation, models.Swift{})

	router := routes.SetupRouter(swiftRepo, validate)
	err = router.Run(":8080")

	if err != nil {
		panic(err)
	}
}
