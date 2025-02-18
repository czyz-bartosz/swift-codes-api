package utils

import (
	"awesomeProject/configs"
	"awesomeProject/dbs"
	"awesomeProject/models"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/uptrace/bun"
	"log"
	"os"
	"strings"
)

func GetFilePath() string {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run import_csv.go <path_to_csv>")
		return ""
	}

	return os.Args[1]
}

func isSwiftCodeOfHeadquarter(swiftCode string) bool {
	return strings.HasSuffix(swiftCode, "XXX")
}

func parseCSVFile(csvFilePath string) ([]models.Bank, error) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %v", csvFilePath, err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	reader := csv.NewReader(file)

	//skip header
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	var banks []models.Bank

	for record, err := reader.Read(); err == nil; record, err = reader.Read() {
		if len(record) < 8 {
			continue
		}

		bank := models.Bank{
			CountryIso2: strings.ToUpper(record[0]),
			SwiftCode:   strings.ToUpper(record[1]),
			Name:        record[3],
			Address:     record[4],
			TownName:    record[5],
			CountryName: strings.ToUpper(record[6]),
		}

		if len(bank.SwiftCode) != 11 {
			continue
		}

		bank.IsHeadquarter = isSwiftCodeOfHeadquarter(bank.SwiftCode)

		banks = append(banks, bank)
	}

	return banks, nil
}

func ImportData(csvFilePath string) error {
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

	ctx := context.Background()

	banks, err := parseCSVFile(csvFilePath)

	if err != nil {
		return err
	}

	err = db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(&banks).Exec(ctx)
		return err
	})

	if err != nil {
		return err
	}

	fmt.Println("Import data successfully")
	return nil
}
