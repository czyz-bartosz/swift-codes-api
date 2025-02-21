package repositories

import (
	"awesomeProject/dbs"
	"awesomeProject/models"
	"context"
	"database/sql"
	"fmt"
)

type SwiftRepoPostgres struct {
	Db dbs.SwiftDb
}

func (swiftRepo SwiftRepoPostgres) GetBySwiftCode(ctx context.Context, swiftCode string) (*models.Swift, error) {
	swift := &models.Swift{}
	err := swiftRepo.Db.NewSelect().Model(swift).Where("swift_code = ?", swiftCode).Scan(ctx)
	return swift, err
}

func (swiftRepo SwiftRepoPostgres) GetBranchesBySwiftCode(ctx context.Context, swiftCode string) ([]models.SwiftMini, error) {
	if len(swiftCode) != 11 {
		return nil, fmt.Errorf("swiftCode must be 11 characters")
	}

	branches := make([]models.SwiftMini, 0)
	query := `
        SELECT address, swifts.bank_name, country_iso2_code, is_headquarter, swift_code 
        FROM swifts 
        WHERE LEFT(swift_code, 8) = ? AND swift_code != ?
    `

	err := swiftRepo.Db.NewRaw(query, swiftCode[:8], swiftCode).Scan(ctx, &branches)

	return branches, err
}

func (swiftRepo SwiftRepoPostgres) GetByCountryIso2Code(ctx context.Context, countryIso2Code string) ([]models.SwiftMini, error) {
	branches := make([]models.SwiftMini, 0)
	query := `
        SELECT address, bank_name, country_iso2_code, is_headquarter, swift_code 
        FROM swifts 
        WHERE swifts.country_iso2_code = ?
    `

	err := swiftRepo.Db.NewRaw(query, countryIso2Code).Scan(ctx, &branches)
	if len(branches) == 0 {
		return nil, sql.ErrNoRows
	}
	return branches, err
}

func (swiftRepo SwiftRepoPostgres) GetCountryNameByIso2Code(ctx context.Context, countryIso2Code string) (string, error) {
	countryName := ""

	query := `
        SELECT swifts.country_name
        FROM swifts
        WHERE swifts.country_iso2_code = ?
        LIMIT 1
    `

	err := swiftRepo.Db.NewRaw(query, countryIso2Code).Scan(ctx, &countryName)
	fmt.Println(countryIso2Code, countryName)

	return countryName, err
}

func (swiftRepo SwiftRepoPostgres) AddSwift(ctx context.Context, swift *models.Swift) error {
	_, err := swiftRepo.Db.NewInsert().Model(swift).Exec(ctx)

	return err
}

func (swiftRepo SwiftRepoPostgres) DeleteSwift(ctx context.Context, swiftCode string) error {
	_, err := swiftRepo.Db.NewDelete().Model(&models.Swift{}).Where("swift_code = ?", swiftCode).Exec(ctx)

	return err
}
