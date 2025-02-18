package models

type BankBranch struct {
	CountryIso2   string `bun:"country_iso2_code" json:"countryISO2"`
	SwiftCode     string `bun:"swift_code," json:"swiftCode"`
	Name          string `bun:"name" json:"bankName"`
	Address       string `bun:"address" json:"address"`
	IsHeadquarter bool   `bun:"is_headquarter" json:"isHeadquarter"`
}
