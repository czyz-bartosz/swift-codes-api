package main

import (
	"awesomeProject/configs"
	"awesomeProject/dbs"
	"awesomeProject/internal/dbimporter/utils"
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"log"
)

func main() {
	csvFilePath := utils.GetFilePath()

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

	ctx := context.Background()

	banks := utils.ParseCSVFile(csvFilePath)

	err := db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(banks).Exec(ctx)
		return err
	})

	if err != nil {
		log.Fatal("Transaction failed: ", err)
		return
	}

	fmt.Println("Done")
}
