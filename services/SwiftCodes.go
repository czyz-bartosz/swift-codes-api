package services

import (
	"awesomeProject/customErrors"
	"awesomeProject/models"
	"awesomeProject/repositories"
	"context"
	"database/sql"
	"errors"
)

func GetBankDetails(ctx context.Context, swiftCode string, bankRepo repositories.BankRepo) (
	bank *models.Bank,
	branches []models.BankBranch,
	err error,
) {
	bank, err = bankRepo.GetBySwiftCode(ctx, swiftCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = customErrors.ErrBankNotFound
			return
		}
		return
	}
	branches, err = bankRepo.GetBranchesBySwiftCode(ctx, swiftCode)
	if err != nil {
		return
	}

	return bank, branches, nil
}
