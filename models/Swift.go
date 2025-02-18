package models

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/uptrace/bun"
	"strings"
)

type Swift struct {
	bun.BaseModel `bun:"table:swifts,alias:s"`

	CountryIso2   string `bun:"country_iso2_code,notnull" json:"countryISO2" validate:"required,iso3166_1_alpha2"`
	SwiftCode     string `bun:"swift_code,pk," json:"swiftCode" validate:"required,len=11,alpha,uppercase"`
	BankName      string `bun:"bank_name,notnull" json:"bankName" validate:"required"`
	Address       string `bun:"address,notnull" json:"address" validate:"required"`
	CountryName   string `bun:"country_name,notnull" json:"countryName" validate:"required"`
	IsHeadquarter bool   `bun:"is_headquarter,notnull" json:"isHeadquarter" validate:"boolean"`
}

func (s *Swift) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	s.SwiftCode = strings.ToUpper(s.SwiftCode)
	s.CountryIso2 = strings.ToUpper(s.CountryIso2)
	s.CountryName = strings.ToUpper(s.CountryName)
	return nil
}

func IsSwiftCodeOfHeadquarter(swiftCode string) bool {
	return strings.HasSuffix(strings.ToUpper(swiftCode), "XXX")
}

func SwiftStructLevelValidation(sl validator.StructLevel) {
	swift := sl.Current().Interface().(Swift)

	if IsSwiftCodeOfHeadquarter(swift.SwiftCode) != swift.IsHeadquarter {
		sl.ReportError(swift.IsHeadquarter, "isHeadquarter", "IsHeadquarter",
			"swiftCode_isHeadquarter_inconsistency", "")
	}
}
