package repositories

import (
	"awesomeProject/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
)

type BankRepoPostgres struct {
	Db *bun.DB
}

func (bankRepo BankRepoPostgres) GetBySwiftCode(ctx context.Context, swiftCode string) (*models.Bank, error) {
	bank := &models.Bank{}
	err := bankRepo.Db.NewSelect().Model(bank).Where("swift_code = ?", swiftCode).Scan(ctx)
	return bank, err
}

func (bankRepo BankRepoPostgres) GetBranchesBySwiftCode(ctx context.Context, swiftCode string) ([]models.BankMini, error) {
	if len(swiftCode) != 11 {
		return nil, fmt.Errorf("swiftCode must be 11 characters")
	}

	branches := make([]models.BankMini, 0)
	query := `
        SELECT address, name, country_iso2_code, is_headquarter, swift_code 
        FROM banks 
        WHERE LEFT(swift_code, 8) = ? AND swift_code != ?
    `

	err := bankRepo.Db.NewRaw(query, swiftCode[:8], swiftCode).Scan(ctx, &branches)
	if len(branches) == 0 {
		return nil, sql.ErrNoRows
	}
	return branches, err
}

func (bankRepo BankRepoPostgres) GetByCountryIso2Code(ctx context.Context, countryIso2Code string) ([]models.BankMini, error) {
	branches := make([]models.BankMini, 0)
	query := `
        SELECT address, name, country_iso2_code, is_headquarter, swift_code 
        FROM banks 
        WHERE banks.country_iso2_code = ?
    `

	err := bankRepo.Db.NewRaw(query, countryIso2Code).Scan(ctx, &branches)
	if len(branches) == 0 {
		return nil, sql.ErrNoRows
	}
	return branches, err
}

func (bankRepo BankRepoPostgres) GetCountryNameByIso2Code(ctx context.Context, countryIso2Code string) (*string, error) {
	countryName := ""

	query := `
        SELECT banks.country_name
        FROM banks 
        WHERE banks.country_iso2_code = ?
        LIMIT 1
    `

	err := bankRepo.Db.NewRaw(query, countryIso2Code).Scan(ctx, &countryName)
	fmt.Println(countryIso2Code, countryName)

	return &countryName, err
}

func (bankRepo BankRepoPostgres) AddBank(ctx context.Context, bank *models.Bank) error {
	_, err := bankRepo.Db.NewInsert().Model(bank).Exec(ctx)

	return err
}
