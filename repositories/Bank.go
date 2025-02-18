package repositories

import (
	"awesomeProject/models"
	"context"
)

type BankRepo interface {
	GetBySwiftCode(context.Context, string) (*models.Bank, error)
	GetBranchesBySwiftCode(context.Context, string) ([]models.BankMini, error)
	GetByCountryIso2Code(context.Context, string) ([]models.BankMini, error)
	GetCountryNameByIso2Code(context.Context, string) (*string, error)
}
