package models

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsSwiftCodeOfHeadquarter(t *testing.T) {
	tests := []struct {
		name      string
		swiftCode string
		want      bool
	}{
		{
			name:      "Headquarter code",
			swiftCode: "ABCDEF12XXX",
			want:      true,
		},
		{
			name:      "Branch code",
			swiftCode: "ABCDEF12345",
			want:      false,
		},
		{
			name:      "Empty code",
			swiftCode: "",
			want:      false,
		},
		{
			name:      "Lowercase headquarter code",
			swiftCode: "abcdef12xxx",
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSwiftCodeOfHeadquarter(tt.swiftCode); got != tt.want {
				t.Errorf("IsSwiftCodeOfHeadquarter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSwiftStructLevelValidation(t *testing.T) {
	validate := validator.New()
	validate.RegisterStructValidation(SwiftStructLevelValidation, Swift{})

	tests := []struct {
		name        string
		swift       Swift
		expectError bool
	}{
		// Poprawne przypadki
		{
			name: "Valid headquarter",
			swift: Swift{
				CountryIso2:   "US",
				SwiftCode:     "ABCDEFDSXXX",
				BankName:      "Bank of Test",
				Address:       "123 Main St",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			expectError: false,
		},
		{
			name: "Valid branch",
			swift: Swift{
				CountryIso2:   "PL",
				SwiftCode:     "ABCDEFHDSAC",
				BankName:      "Bank of Poland",
				Address:       "456 Warsaw St",
				CountryName:   "Poland",
				IsHeadquarter: false,
			},
			expectError: false,
		},

		// Niepoprawne przypadki
		{
			name: "Invalid headquarter - code ends with XXX but IsHeadquarter is false",
			swift: Swift{
				CountryIso2:   "US",
				SwiftCode:     "ABCDEF12XXX",
				BankName:      "Bank of Test",
				Address:       "123 Main St",
				CountryName:   "United States",
				IsHeadquarter: false,
			},
			expectError: true,
		},
		{
			name: "Invalid branch - code does not end with XXX but IsHeadquarter is true",
			swift: Swift{
				CountryIso2:   "PL",
				SwiftCode:     "ABCDEF12345",
				BankName:      "Bank of Poland",
				Address:       "456 Warsaw St",
				CountryName:   "Poland",
				IsHeadquarter: true,
			},
			expectError: true,
		},
		{
			name: "Invalid CountryIso2 - not ISO 3166-1 alpha-2",
			swift: Swift{
				CountryIso2:   "USA", // Niepoprawny format
				SwiftCode:     "ABCDEF12XXX",
				BankName:      "Bank of Test",
				Address:       "123 Main St",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			expectError: true,
		},
		{
			name: "Invalid SwiftCode - not 11 characters",
			swift: Swift{
				CountryIso2:   "US",
				SwiftCode:     "ABCDEF12", // Za krótki
				BankName:      "Bank of Test",
				Address:       "123 Main St",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			expectError: true,
		},
		{
			name: "Invalid SwiftCode - contains special characters",
			swift: Swift{
				CountryIso2:   "US",
				SwiftCode:     "ABCDEF12@#$", // Znaki specjalne
				BankName:      "Bank of Test",
				Address:       "123 Main St",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			expectError: true,
		},
		{
			name: "Invalid SwiftCode - contains lowercase letters",
			swift: Swift{
				CountryIso2:   "US",
				SwiftCode:     "abcdef12xxx", // Małe litery
				BankName:      "Bank of Test",
				Address:       "123 Main St",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			expectError: true,
		},
		{
			name: "Missing required fields",
			swift: Swift{
				CountryIso2:   "", // Brak wymaganego pola
				SwiftCode:     "", // Brak wymaganego pola
				BankName:      "", // Brak wymaganego pola
				Address:       "", // Brak wymaganego pola
				CountryName:   "", // Brak wymaganego pola
				IsHeadquarter: true,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.swift)
			if tt.expectError {
				assert.Error(t, err, "Expected an error for case: %s", tt.name)
			} else {
				assert.NoError(t, err, "Expected no error for case: %s", tt.name)
			}
		})
	}
}

func TestSwift_BeforeAppendModel(t *testing.T) {
	tests := []struct {
		name         string
		swift        Swift
		expectedCode string
		expectedIso2 string
		expectedName string
	}{
		{
			name: "Convert to uppercase",
			swift: Swift{
				SwiftCode:   "abcdef12xxx",
				CountryIso2: "us",
				CountryName: "united states",
			},
			expectedCode: "ABCDEF12XXX",
			expectedIso2: "US",
			expectedName: "UNITED STATES",
		},
		{
			name: "Already uppercase",
			swift: Swift{
				SwiftCode:   "ABCDEF12XXX",
				CountryIso2: "US",
				CountryName: "UNITED STATES",
			},
			expectedCode: "ABCDEF12XXX",
			expectedIso2: "US",
			expectedName: "UNITED STATES",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.swift.BeforeAppendModel(context.Background(), nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, tt.swift.SwiftCode)
			assert.Equal(t, tt.expectedIso2, tt.swift.CountryIso2)
			assert.Equal(t, tt.expectedName, tt.swift.CountryName)
		})
	}
}
