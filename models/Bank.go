package models

import (
	"context"
	"github.com/uptrace/bun"
	"strings"
)

type Bank struct {
	bun.BaseModel `bun:"table:banks,alias:b"`

	CountryIso2   string `bun:"country_iso2_code,notnull" json:"countryISO2"`
	SwiftCode     string `bun:"swift_code,pk," json:"swiftCode"`
	Name          string `bun:"name,notnull" json:"bankName"`
	Address       string `bun:"address,notnull" json:"address"`
	TownName      string `bun:"town_name,notnull" json:"townName,omitempty"`
	CountryName   string `bun:"country_name,notnull" json:"countryName"`
	IsHeadquarter bool   `bun:"is_headquarter,notnull" json:"isHeadquarter"`
}

func (b *Bank) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	b.SwiftCode = strings.ToUpper(b.SwiftCode)
	b.CountryIso2 = strings.ToUpper(b.CountryIso2)
	b.CountryName = strings.ToUpper(b.CountryName)
	return nil
}

func IsSwiftCodeOfHeadquarter(swiftCode string) bool {
	return strings.HasSuffix(swiftCode, "XXX")
}
