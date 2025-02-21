package services

import (
	"awesomeProject/customErrors"
	"awesomeProject/mocks"
	"awesomeProject/models"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetSwiftDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := &SwiftServiceDefault{}
	mockSwiftRepo := mocks.NewMockSwiftRepo(ctrl)
	ctx := context.Background()

	tests := []struct {
		name         string
		swiftCode    string
		mockSetup    func()
		wantSwift    *models.Swift
		wantBranches []models.SwiftMini
		wantErr      error
	}{
		{
			name:      "Success - Get Swift Details",
			swiftCode: "ABCDEFGHXXX",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "ABCDEFGHXXX").Return(&models.Swift{
					SwiftCode:     "ABCDEFGHXXX",
					BankName:      "Test Bank",
					Address:       "123 Test St",
					CountryIso2:   "US",
					CountryName:   "United States",
					IsHeadquarter: true,
				}, nil)
				mockSwiftRepo.EXPECT().GetBranchesBySwiftCode(ctx, "ABCDEFGHXXX").Return([]models.SwiftMini{
					{
						SwiftCode:     "ABCDEFGH001",
						BankName:      "Test Bank Branch 1",
						Address:       "456 Branch St",
						CountryIso2:   "US",
						IsHeadquarter: false,
					},
				}, nil)
			},
			wantSwift: &models.Swift{
				SwiftCode:     "ABCDEFGHXXX",
				BankName:      "Test Bank",
				Address:       "123 Test St",
				CountryIso2:   "US",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			wantBranches: []models.SwiftMini{
				{
					SwiftCode:     "ABCDEFGH001",
					BankName:      "Test Bank Branch 1",
					Address:       "456 Branch St",
					CountryIso2:   "US",
					IsHeadquarter: false,
				},
			},
			wantErr: nil,
		},
		{
			name:      "Error - Swift Not Found",
			swiftCode: "INVALIDCODE",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "INVALIDCODE").Return(nil, sql.ErrNoRows)
			},
			wantSwift:    nil,
			wantBranches: nil,
			wantErr:      customErrors.ErrSwiftNotFound,
		},
		{
			name:      "Success - Swift Code Not Headquarter (No Branches)",
			swiftCode: "ABCDEFGH001",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "ABCDEFGH001").Return(&models.Swift{
					SwiftCode:     "ABCDEFGH001",
					BankName:      "Test Bank Branch",
					Address:       "456 Branch St",
					CountryIso2:   "US",
					CountryName:   "United States",
					IsHeadquarter: false,
				}, nil)
			},
			wantSwift: &models.Swift{
				SwiftCode:     "ABCDEFGH001",
				BankName:      "Test Bank Branch",
				Address:       "456 Branch St",
				CountryIso2:   "US",
				CountryName:   "United States",
				IsHeadquarter: false,
			},
			wantBranches: nil,
			wantErr:      nil,
		},
		{
			name:      "Error - unknown error",
			swiftCode: "",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "").Return(nil, errors.New("db error"))
			},
			wantSwift:    nil,
			wantBranches: nil,
			wantErr:      errors.New("db error"),
		},
		{
			name:      "Success - Get Swift Details",
			swiftCode: "ABCDEFGHXXX",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "ABCDEFGHXXX").Return(&models.Swift{
					SwiftCode:     "ABCDEFGHXXX",
					BankName:      "Test Bank",
					Address:       "123 Test St",
					CountryIso2:   "US",
					CountryName:   "United States",
					IsHeadquarter: true,
				}, nil)
				mockSwiftRepo.EXPECT().GetBranchesBySwiftCode(ctx, "ABCDEFGHXXX").Return(nil, errors.New("db error"))
			},
			wantSwift: &models.Swift{
				SwiftCode:     "ABCDEFGHXXX",
				BankName:      "Test Bank",
				Address:       "123 Test St",
				CountryIso2:   "US",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			wantBranches: nil,
			wantErr:      errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotSwift, gotBranches, err := service.GetSwiftDetails(ctx, tt.swiftCode, mockSwiftRepo)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantSwift, gotSwift)
			assert.Equal(t, tt.wantBranches, gotBranches)
		})
	}
}

func TestGetSwiftsDetailsByCountryIso2Code(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := &SwiftServiceDefault{}
	mockSwiftRepo := mocks.NewMockSwiftRepo(ctrl)
	ctx := context.Background()

	tests := []struct {
		name            string
		countryIso2Code string
		mockSetup       func()
		wantCountryName string
		wantSwifts      []models.SwiftMini
		wantErr         error
	}{
		{
			name:            "Success - Get Swifts by Country ISO2 Code",
			countryIso2Code: "US",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetByCountryIso2Code(ctx, "US").Return([]models.SwiftMini{
					{
						SwiftCode:     "ABCDEFGHXXX",
						BankName:      "Test Bank",
						Address:       "123 Test St",
						CountryIso2:   "US",
						IsHeadquarter: true,
					},
				}, nil)
				mockSwiftRepo.EXPECT().GetCountryNameByIso2Code(ctx, "US").Return("United States", nil)
			},
			wantCountryName: "United States",
			wantSwifts: []models.SwiftMini{
				{
					SwiftCode:     "ABCDEFGHXXX",
					BankName:      "Test Bank",
					Address:       "123 Test St",
					CountryIso2:   "US",
					IsHeadquarter: true,
				},
			},
			wantErr: nil,
		},
		{
			name:            "Error - Country Not Found",
			countryIso2Code: "XX",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetByCountryIso2Code(ctx, "XX").Return(nil, sql.ErrNoRows)
			},
			wantCountryName: "",
			wantSwifts:      nil,
			wantErr:         customErrors.ErrSwiftNotFound,
		},
		{
			name:            "Error - unknown error",
			countryIso2Code: "XX",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetByCountryIso2Code(ctx, "XX").Return(nil, errors.New("db error"))
			},
			wantCountryName: "",
			wantSwifts:      nil,
			wantErr:         errors.New("db error"),
		},
		{
			name:            "Error - unknown error",
			countryIso2Code: "XX",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetByCountryIso2Code(ctx, "XX").Return(nil, nil)
				mockSwiftRepo.EXPECT().GetCountryNameByIso2Code(ctx, "XX").Return("", errors.New("db error"))
			},
			wantCountryName: "",
			wantSwifts:      nil,
			wantErr:         errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotCountryName, gotSwifts, err := service.GetSwiftsDetailsByCountryIso2Code(ctx, tt.countryIso2Code, mockSwiftRepo)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantCountryName, gotCountryName)
			assert.Equal(t, tt.wantSwifts, gotSwifts)
		})
	}
}

func TestAddSwift(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := &SwiftServiceDefault{}
	mockSwiftRepo := mocks.NewMockSwiftRepo(ctrl)
	mockValidator := mocks.NewMockSwiftValidator(ctrl)
	ctx := context.Background()

	tests := []struct {
		name      string
		swift     *models.Swift
		mockSetup func()
		wantErr   error
	}{
		{
			name: "Success - Add Swift",
			swift: &models.Swift{
				SwiftCode:     "aBCDEFGhXXX",
				BankName:      "Test Bank",
				Address:       "123 Test St",
				CountryIso2:   "US",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			mockSetup: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "ABCDEFGHXXX").Return(nil, sql.ErrNoRows)
				mockSwiftRepo.EXPECT().AddSwift(ctx, gomock.Any()).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "Error - Swift Code Already Exists",
			swift: &models.Swift{
				SwiftCode:     "ABCDEFGHXXX",
				BankName:      "Test Bank",
				Address:       "123 Test St",
				CountryIso2:   "US",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			mockSetup: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "ABCDEFGHXXX").Return(&models.Swift{}, nil)
			},
			wantErr: customErrors.ErrSwiftCodeAlreadyExists,
		},
		{
			name: "Error - Swift Code Already Exists",
			swift: &models.Swift{
				SwiftCode:     "ABCDEFGHXXX",
				BankName:      "Test Bank",
				Address:       "123 Test St",
				CountryIso2:   "US",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			mockSetup: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "ABCDEFGHXXX").Return(&models.Swift{}, nil)
			},
			wantErr: customErrors.ErrSwiftCodeAlreadyExists,
		},
		{
			name: "Error - validation error",
			swift: &models.Swift{
				SwiftCode:     "ABCDEFGHXXX",
				BankName:      "Test Bank",
				Address:       "123 Test St",
				CountryIso2:   "US",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			mockSetup: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(errors.New("validation error"))
				mockSwiftRepo.EXPECT().AddSwift(ctx, gomock.Any()).Times(0)
			},
			wantErr: errors.New("validation error"),
		},
		{
			name: "Error - unknown error",
			swift: &models.Swift{
				SwiftCode:     "ABCDEFGHXXX",
				BankName:      "Test Bank",
				Address:       "123 Test St",
				CountryIso2:   "US",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			mockSetup: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "ABCDEFGHXXX").Return(nil, errors.New("db error"))
				mockSwiftRepo.EXPECT().AddSwift(ctx, gomock.Any()).Times(0)
			},
			wantErr: errors.New("db error"),
		},
		{
			name: "Error - unknown error",
			swift: &models.Swift{
				SwiftCode:     "ABCDEFGHXXX",
				BankName:      "Test Bank",
				Address:       "123 Test St",
				CountryIso2:   "US",
				CountryName:   "United States",
				IsHeadquarter: true,
			},
			mockSetup: func() {
				mockValidator.EXPECT().Struct(gomock.Any()).Return(nil)
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "ABCDEFGHXXX").Return(nil, sql.ErrNoRows)
				mockSwiftRepo.EXPECT().AddSwift(ctx, gomock.Any()).Return(errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.AddSwift(ctx, tt.swift, mockSwiftRepo, mockValidator)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeleteSwift(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := &SwiftServiceDefault{}
	mockSwiftRepo := mocks.NewMockSwiftRepo(ctrl)
	ctx := context.Background()

	tests := []struct {
		name      string
		swiftCode string
		mockSetup func()
		wantErr   error
	}{
		{
			name:      "Success - Delete Swift",
			swiftCode: "ABcDEFGHXXX",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "ABCDEFGHXXX").Return(&models.Swift{}, nil)
				mockSwiftRepo.EXPECT().DeleteSwift(ctx, "ABCDEFGHXXX").Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "Error - Swift Not Found",
			swiftCode: "INVALIDCODE",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "INVALIDCODE").Return(nil, sql.ErrNoRows)
			},
			wantErr: customErrors.ErrSwiftNotFound,
		},
		{
			name:      "Error - unknown error",
			swiftCode: "INVALIDCODE",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "INVALIDCODE").Return(nil, errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
		{
			name:      "Error - unknown error",
			swiftCode: "INVALIDCODE",
			mockSetup: func() {
				mockSwiftRepo.EXPECT().GetBySwiftCode(ctx, "INVALIDCODE").Return(nil, nil)
				mockSwiftRepo.EXPECT().DeleteSwift(ctx, "INVALIDCODE").Return(errors.New("db error"))
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := service.DeleteSwift(ctx, tt.swiftCode, mockSwiftRepo)
			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
