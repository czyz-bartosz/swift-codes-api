package models

import (
	"context"
	"github.com/uptrace/bun"
	"strings"
)

type Bank struct {
	bun.BaseModel `bun:"table:banks,alias:b"`

	CountryIso2   string `bun:"country_iso2_code,notnull" json:"country_iso2_code"`
	SwiftCode     string `bun:"swift_code,pk," json:"swift_code"`
	Name          string `bun:"name,notnull" json:"name"`
	Address       string `bun:"address,notnull" json:"address"`
	TownName      string `bun:"town_name,notnull" json:"town_name"`
	CountryName   string `bun:"country_name,notnull" json:"country_name"`
	IsHeadquarter bool   `bun:"is_headquarter,notnull" json:"is_headquarter"`
}

func (b *Bank) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	b.SwiftCode = strings.ToUpper(b.SwiftCode)
	b.CountryIso2 = strings.ToUpper(b.CountryIso2)
	b.CountryName = strings.ToUpper(b.CountryName)
	return nil
}
