package services

import (
	"awesomeProject/customErrors"
	"awesomeProject/models"
	"awesomeProject/repositories"
	"context"
	"database/sql"
	"errors"
	"github.com/go-playground/validator/v10"
	"strings"
)

func GetSwiftDetails(ctx context.Context, swiftCode string, bankRepo repositories.SwiftRepo) (
	swift *models.Swift,
	branches []models.SwiftMini,
	err error,
) {
	swift, err = bankRepo.GetBySwiftCode(ctx, swiftCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = customErrors.ErrSwiftNotFound
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
	countryName string,
	swifts []models.SwiftMini,
	err error,
) {
	swifts, err = bankRepo.GetByCountryIso2Code(ctx, countryIso2Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = customErrors.ErrSwiftNotFound
		}
		return
	}

	countryName, err = bankRepo.GetCountryNameByIso2Code(ctx, countryIso2Code)
	if err != nil {
		return
	}

	return
}

func AddSwift(ctx context.Context, swift *models.Swift, bankRepo repositories.SwiftRepo, validate *validator.Validate) error {
	err := validate.Struct(swift)
	if err != nil {
		return err
	}
	swift.SwiftCode = strings.ToUpper(swift.SwiftCode)
	_, err = bankRepo.GetBySwiftCode(ctx, swift.SwiftCode)

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

func DeleteSwift(ctx context.Context, swiftCode string, bankRepo repositories.SwiftRepo) error {
	swiftCode = strings.ToUpper(swiftCode)
	_, err := bankRepo.GetBySwiftCode(ctx, swiftCode)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customErrors.ErrSwiftNotFound
		}
		return err
	}

	err = bankRepo.DeleteSwift(ctx, swiftCode)

	return err
}
