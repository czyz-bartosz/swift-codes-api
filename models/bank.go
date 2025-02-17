package models

import (
	"context"
	"github.com/uptrace/bun"
	"strings"
)

type Bank struct {
	bun.BaseModel `bun:"table:banks,alias:b"`

	ID            int    `bun:"id,pk,autoincrement" json:"id"`
	CountryIso2   string `bun:"country_iso2_code,notnull" json:"country_iso2_code"`
	SwiftCode     string `bun:"swift_code,notnull,unique" json:"swift_code"`
	Name          string `bun:"name,notnull" json:"name"`
	Address       string `bun:"address,notnull" json:"address"`
	TownName      string `bun:"town_name,notnull" json:"town_name"`
	CountryName   string `bun:"country_name,notnull" json:"country_name"`
	HeadquarterID *int   `bun:"headquarter_id" json:"headquarter_id,omitempty"`

	// Self-referential relationship
	Headquarter *Bank `bun:"rel:belongs-to,join:headquarter_id=id"`
}

func (b *Bank) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	b.SwiftCode = strings.ToUpper(b.SwiftCode)
	b.CountryIso2 = strings.ToUpper(b.CountryIso2)
	b.CountryName = strings.ToUpper(b.CountryName)
	return nil
}
