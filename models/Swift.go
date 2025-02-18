package models

import (
	"context"
	"github.com/uptrace/bun"
	"strings"
)

type Swift struct {
	bun.BaseModel `bun:"table:swifts,alias:s"`

	CountryIso2   string `bun:"country_iso2_code,notnull" json:"countryISO2"`
	SwiftCode     string `bun:"swift_code,pk," json:"swiftCode"`
	BankName      string `bun:"bank_name,notnull" json:"bankName"`
	Address       string `bun:"address,notnull" json:"address"`
	TownName      string `bun:"town_name,notnull" json:"townName,omitempty"`
	CountryName   string `bun:"country_name,notnull" json:"countryName"`
	IsHeadquarter bool   `bun:"is_headquarter,notnull" json:"isHeadquarter"`
}

func (s *Swift) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	s.SwiftCode = strings.ToUpper(s.SwiftCode)
	s.CountryIso2 = strings.ToUpper(s.CountryIso2)
	s.CountryName = strings.ToUpper(s.CountryName)
	return nil
}

func IsSwiftCodeOfHeadquarter(swiftCode string) bool {
	return strings.HasSuffix(swiftCode, "XXX")
}
