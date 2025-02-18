package repositories

import (
	"awesomeProject/models"
	"context"
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

func (bankRepo BankRepoPostgres) GetBranchesBySwiftCode(ctx context.Context, swiftCode string) ([]models.BankBranch, error) {
	if len(swiftCode) != 11 {
		return nil, fmt.Errorf("swiftCode must be 11 characters")
	}

	branches := make([]models.BankBranch, 0)
	query := `
        SELECT address, name, country_iso2_code, is_headquarter, swift_code 
        FROM banks 
        WHERE LEFT(swift_code, 8) = ? AND swift_code != ?
    `

	err := bankRepo.Db.NewRaw(query, swiftCode[:8], swiftCode).Scan(ctx, &branches)
	return branches, err
}
