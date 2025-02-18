package migrations

import (
	"awesomeProject/models"
	"context"
	"fmt"
	"github.com/uptrace/bun"
)

func Migrate(db *bun.DB) error {
	return db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		if _, err := tx.NewCreateTable().
			IfNotExists().
			Model((*models.Swift)(nil)).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}

		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*models.Swift)(nil)).
			Index("idx_swift_code_prefix").
			ColumnExpr("LEFT(swift_code, 8)").
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create index on swift_code prefix: %w", err)
		}

		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*models.Swift)(nil)).
			Index("idx_country_iso2_code").
			Column("country_iso2_code").
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create index on country_iso2_code: %w", err)
		}

		fmt.Println("Migrations done")
		return nil
	})
}
