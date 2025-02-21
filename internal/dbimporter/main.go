package main

import (
	"awesomeProject/configs"
	"awesomeProject/dbs"
	"awesomeProject/internal/dbimporter/utils"
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
			fmt.Println("Error closing db")
		}
	}(db)
	csvFilePath := utils.GetFilePath()
	err := utils.ImportData(csvFilePath, db)
	if err != nil {
		panic(err)
	}
}
