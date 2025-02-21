package services

import (
	"awesomeProject/customErrors"
	"awesomeProject/models"
	"awesomeProject/repositories"
	"context"
	"database/sql"
	"errors"
	"strings"
)

type SwiftService interface {
	GetSwiftDetails(ctx context.Context, swiftCode string, swiftRepo repositories.SwiftRepo) (
		swift *models.Swift,
		branches []models.SwiftMini,
		err error,
	)
	GetSwiftsDetailsByCountryIso2Code(ctx context.Context, countryIso2Code string, swiftRepo repositories.SwiftRepo) (
		countryName string,
		swifts []models.SwiftMini,
		err error,
	)
	AddSwift(ctx context.Context, swift *models.Swift, swiftRepo repositories.SwiftRepo, validate models.SwiftValidator) error
	DeleteSwift(ctx context.Context, swiftCode string, swiftRepo repositories.SwiftRepo) error
}

type SwiftServiceDefault struct{}

func (s *SwiftServiceDefault) GetSwiftDetails(ctx context.Context, swiftCode string, swiftRepo repositories.SwiftRepo) (
	swift *models.Swift,
	branches []models.SwiftMini,
	err error,
) {
	swift, err = swiftRepo.GetBySwiftCode(ctx, swiftCode)
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

	branches, err = swiftRepo.GetBranchesBySwiftCode(ctx, swiftCode)
	if err != nil {
		return
	}

	return swift, branches, nil
}

func (s *SwiftServiceDefault) GetSwiftsDetailsByCountryIso2Code(ctx context.Context, countryIso2Code string, swiftRepo repositories.SwiftRepo) (
	countryName string,
	swifts []models.SwiftMini,
	err error,
) {
	swifts, err = swiftRepo.GetByCountryIso2Code(ctx, countryIso2Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = customErrors.ErrSwiftNotFound
		}
		return
	}

	countryName, err = swiftRepo.GetCountryNameByIso2Code(ctx, countryIso2Code)
	if err != nil {
		return
	}

	return
}

func (s *SwiftServiceDefault) AddSwift(ctx context.Context, swift *models.Swift, swiftRepo repositories.SwiftRepo, validate models.SwiftValidator) error {
	err := validate.Struct(swift)
	if err != nil {
		return err
	}
	swift.SwiftCode = strings.ToUpper(swift.SwiftCode)
	_, err = swiftRepo.GetBySwiftCode(ctx, swift.SwiftCode)

	if err == nil {
		return customErrors.ErrSwiftCodeAlreadyExists
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	err = swiftRepo.AddSwift(ctx, swift)
	if err != nil {
		return err
	}

	return nil
}

func (s *SwiftServiceDefault) DeleteSwift(ctx context.Context, swiftCode string, swiftRepo repositories.SwiftRepo) error {
	swiftCode = strings.ToUpper(swiftCode)
	_, err := swiftRepo.GetBySwiftCode(ctx, swiftCode)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return customErrors.ErrSwiftNotFound
		}
		return err
	}

	err = swiftRepo.DeleteSwift(ctx, swiftCode)

	return err
}
