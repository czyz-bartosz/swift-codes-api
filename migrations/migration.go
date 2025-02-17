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
			Model((*models.Bank)(nil)).
			ForeignKey(`("headquarter_id") REFERENCES "banks" ("id") ON DELETE CASCADE`).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}

		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*models.Bank)(nil)).
			Index("idx_swift_code_full").
			Column("swift_code").
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create index on swift_code: %w", err)
		}

		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*models.Bank)(nil)).
			Index("idx_swift_code_prefix").
			ColumnExpr("SUBSTRING(swift_code FROM 1 FOR 8)").
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create index on swift_code prefix: %w", err)
		}

		// Dodanie indeksu na country_iso2_code
		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*models.Bank)(nil)).
			Index("idx_country_iso2_code").
			Column("country_iso2_code").
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create index on country_iso2_code: %w", err)
		}

		if _, err := tx.NewCreateIndex().
			IfNotExists().
			Model((*models.Bank)(nil)).
			Index("idx_headquarter_id").
			Column("headquarter_id").
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to create index on headquarter_id: %w", err)
		}
		fmt.Println("Migrations done")
		return nil
	})
}
