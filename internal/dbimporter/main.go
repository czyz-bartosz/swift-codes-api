package main

import (
	"awesomeProject/internal/dbimporter/utils"
)

func main() {
	csvFilePath := utils.GetFilePath()
	err := utils.ImportData(csvFilePath)
	if err != nil {
		panic(err)
	}
}
