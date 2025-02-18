package models

type SwiftMini struct {
	CountryIso2   string `bun:"country_iso2_code" json:"countryISO2"`
	SwiftCode     string `bun:"swift_code," json:"swiftCode"`
	BankName      string `bun:"bank_name" json:"bankName"`
	Address       string `bun:"address" json:"address"`
	IsHeadquarter bool   `bun:"is_headquarter" json:"isHeadquarter"`
}
