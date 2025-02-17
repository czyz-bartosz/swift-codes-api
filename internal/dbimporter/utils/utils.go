package utils

import (
	"awesomeProject/models"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func GetFilePath() string {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run import_csv.go <path_to_csv>")
		os.Exit(1)
	}

	return os.Args[1]
}

func isSwiftCodeOfHeadquarter(swiftCode string) bool {
	return strings.HasSuffix(swiftCode, "XXX")
}

func ParseCSVFile(csvFilePath string) *[]*models.Bank {
	file, err := os.Open(csvFilePath)
	if err != nil {
		panic(fmt.Sprintf("could not open file %s: %v", csvFilePath, err))
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// Read the CSV file
	reader := csv.NewReader(file)

	//skip header
	_, err = reader.Read()
	if err != nil {
		panic(err)
	}

	var banks []*models.Bank

	for record, err := reader.Read(); err == nil; record, err = reader.Read() {
		if len(record) != 8 {
			panic(fmt.Sprintf("expected 8 fields, got %d", len(record)))
		}

		bank := &models.Bank{
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

	return &banks
}
