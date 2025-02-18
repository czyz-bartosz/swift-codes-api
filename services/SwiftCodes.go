package services

import (
	"awesomeProject/customErrors"
	"awesomeProject/models"
	"awesomeProject/repositories"
	"context"
	"database/sql"
	"errors"
)

func GetSwiftDetails(ctx context.Context, swiftCode string, bankRepo repositories.SwiftRepo) (
	swift *models.Swift,
	branches []models.SwiftMini,
	err error,
) {
	swift, err = bankRepo.GetBySwiftCode(ctx, swiftCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = customErrors.ErrBankNotFound
			return
		}
		return
	}
	if !models.IsSwiftCodeOfHeadquarter(swiftCode) {
		return swift, nil, nil
	}

	branches, err = bankRepo.GetBranchesBySwiftCode(ctx, swiftCode)
	if err != nil {
		return
	}

	return swift, branches, nil
}

func GetSwiftsDetailsByCountryIso2Code(ctx context.Context, countryIso2Code string, bankRepo repositories.SwiftRepo) (
	countryName *string,
	swifts []models.SwiftMini,
	err error,
) {
	swifts, err = bankRepo.GetByCountryIso2Code(ctx, countryIso2Code)
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

func AddSwift(ctx context.Context, swift *models.Swift, bankRepo repositories.SwiftRepo) error {
	_, err := bankRepo.GetBySwiftCode(ctx, swift.SwiftCode)

	if err == nil {
		return customErrors.ErrSwiftCodeAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = bankRepo.AddSwift(ctx, swift)
	if err != nil {
		return err
	}

	return nil
}
