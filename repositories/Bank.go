package repositories

import (
	"awesomeProject/models"
	"context"
)

type BankRepo interface {
	GetBySwiftCode(context.Context, string) (*models.Bank, error)
	GetBranchesBySwiftCode(context.Context, string) ([]models.BankBranch, error)
}
