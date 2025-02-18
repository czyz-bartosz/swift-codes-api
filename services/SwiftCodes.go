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
	branches []models.BankMini,
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
	if !models.IsSwiftCodeOfHeadquarter(swiftCode) {
		return bank, nil, nil
	}

	branches, err = bankRepo.GetBranchesBySwiftCode(ctx, swiftCode)
	if err != nil {
		return
	}

	return bank, branches, nil
}

func GetBanksDetailsByCountryIso2Code(ctx context.Context, countryIso2Code string, bankRepo repositories.BankRepo) (
	countryName *string,
	banks []models.BankMini,
	err error,
) {
	banks, err = bankRepo.GetByCountryIso2Code(ctx, countryIso2Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = customErrors.ErrBankNotFound
		}
		return
	}

	countryName, err = bankRepo.GetCountryNameByIso2Code(ctx, countryIso2Code)
	if err != nil {
		return
	}

	return
}

func AddBank(ctx context.Context, bank *models.Bank, bankRepo repositories.BankRepo) error {
	_, err := bankRepo.GetBySwiftCode(ctx, bank.SwiftCode)

	if err == nil {
		return customErrors.ErrSwiftCodeAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = bankRepo.AddBank(ctx, bank)
	if err != nil {
		return err
	}

	return nil
}
