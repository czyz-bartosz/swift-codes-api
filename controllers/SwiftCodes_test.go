package controllers

import (
	"awesomeProject/customErrors"
	"awesomeProject/mocks"
	"awesomeProject/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestController_GetSwiftDetails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSwiftRepo := mocks.NewMockSwiftRepo(ctrl)
	mockSwiftService := mocks.NewMockSwiftService(ctrl)

	controller := Controller{
		SwiftRepo:    mockSwiftRepo,
		SwiftService: mockSwiftService,
	}

	tests := []struct {
		name           string
		swiftCode      string
		mockSetup      func()
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:      "Success - Headquarter",
			swiftCode: "ABCDEF12XXX",
			mockSetup: func() {
				mockSwiftService.EXPECT().GetSwiftDetails(gomock.Any(), "ABCDEF12XXX", mockSwiftRepo).Return(
					&models.Swift{
						Address:       "123 Main St",
						BankName:      "Bank of Test",
						CountryIso2:   "US",
						CountryName:   "United States",
						IsHeadquarter: true,
						SwiftCode:     "ABCDEF12XXX",
					},
					[]models.SwiftMini{
						{
							SwiftCode:     "ABCDEF12346",
							BankName:      "Branch 1",
							CountryIso2:   "US",
							IsHeadquarter: false,
							Address:       "456 Branch St",
						},
						{
							SwiftCode:     "ABCDEF12347",
							BankName:      "Branch 2",
							CountryIso2:   "US",
							IsHeadquarter: false,
							Address:       "789 Branch St",
						},
					},
					nil,
				)
			},
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"address":       "123 Main St",
				"bankName":      "Bank of Test",
				"countryISO2":   "US",
				"countryName":   "United States",
				"isHeadquarter": true,
				"swiftCode":     "ABCDEF12XXX",
				"branches": []interface{}{
					map[string]interface{}{
						"swiftCode":     "ABCDEF12346",
						"bankName":      "Branch 1",
						"countryISO2":   "US",
						"isHeadquarter": false,
						"address":       "456 Branch St",
					},
					map[string]interface{}{
						"swiftCode":     "ABCDEF12347",
						"bankName":      "Branch 2",
						"countryISO2":   "US",
						"isHeadquarter": false,
						"address":       "789 Branch St",
					},
				},
			},
		},
		{
			name:      "Success - Branch",
			swiftCode: "ABCDEF12346",
			mockSetup: func() {
				mockSwiftService.EXPECT().GetSwiftDetails(gomock.Any(), "ABCDEF12346", mockSwiftRepo).Return(
					&models.Swift{
						Address:       "456 Branch St",
						BankName:      "Bank of Test Branch",
						CountryIso2:   "US",
						CountryName:   "United States",
						IsHeadquarter: false,
						SwiftCode:     "ABCDEF12346",
					},
					nil,
					nil,
				)
			},
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"address":       "456 Branch St",
				"bankName":      "Bank of Test Branch",
				"countryISO2":   "US",
				"countryName":   "United States",
				"isHeadquarter": false,
				"swiftCode":     "ABCDEF12346",
			},
		},
		{
			name:      "Error - Swift not found",
			swiftCode: "INVALIDCODE",
			mockSetup: func() {
				mockSwiftService.EXPECT().GetSwiftDetails(gomock.Any(), "INVALIDCODE", mockSwiftRepo).Return(
					nil,
					nil,
					customErrors.ErrSwiftNotFound,
				)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"message": "Swift not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/swift/"+tt.swiftCode, nil)
			c.Params = gin.Params{gin.Param{Key: "swiftCode", Value: tt.swiftCode}}

			controller.GetSwiftDetails(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var responseBody gin.H
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func TestController_GetSwiftsDetailsByCountryIso2Code(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSwiftRepo := mocks.NewMockSwiftRepo(ctrl)
	mockSwiftService := mocks.NewMockSwiftService(ctrl)

	controller := Controller{
		SwiftRepo:    mockSwiftRepo,
		SwiftService: mockSwiftService,
	}

	tests := []struct {
		name           string
		countryIso2    string
		mockSetup      func()
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name:        "Success",
			countryIso2: "US",
			mockSetup: func() {
				mockSwiftService.EXPECT().GetSwiftsDetailsByCountryIso2Code(gomock.Any(), "US", mockSwiftRepo).Return(
					"United States",
					[]models.SwiftMini{
						{
							SwiftCode:     "ABCDEF12XXX",
							BankName:      "Bank of Test",
							CountryIso2:   "US",
							IsHeadquarter: true,
							Address:       "123 Main St",
						},
						{
							SwiftCode:     "ABCDEF12346",
							BankName:      "Bank of Test Branch",
							CountryIso2:   "US",
							IsHeadquarter: false,
							Address:       "456 Branch St",
						},
					},
					nil,
				)
			},
			expectedStatus: http.StatusOK,
			expectedBody: gin.H{
				"countryISO2": "US",
				"countryName": "United States",
				"swiftCodes": []interface{}{
					map[string]interface{}{
						"swiftCode":     "ABCDEF12XXX",
						"bankName":      "Bank of Test",
						"countryISO2":   "US",
						"isHeadquarter": true,
						"address":       "123 Main St",
					},
					map[string]interface{}{
						"swiftCode":     "ABCDEF12346",
						"bankName":      "Bank of Test Branch",
						"countryISO2":   "US",
						"isHeadquarter": false,
						"address":       "456 Branch St",
					},
				},
			},
		},
		{
			name:        "Error - Swift not found",
			countryIso2: "INVALID",
			mockSetup: func() {
				mockSwiftService.EXPECT().GetSwiftsDetailsByCountryIso2Code(gomock.Any(), "INVALID", mockSwiftRepo).Return(
					"",
					nil,
					customErrors.ErrSwiftNotFound,
				)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"message": "Swift not found"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/swift/country/"+tt.countryIso2, nil)
			c.Params = gin.Params{gin.Param{Key: "countryIso2Code", Value: tt.countryIso2}}

			controller.GetSwiftsDetailsByCountryIso2Code(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var responseBody gin.H
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func TestController_AddSwift(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSwiftRepo := mocks.NewMockSwiftRepo(ctrl)
	mockSwiftService := mocks.NewMockSwiftService(ctrl)
	mockValidator := mocks.NewMockSwiftValidator(ctrl)

	controller := Controller{
		SwiftRepo:    mockSwiftRepo,
		SwiftService: mockSwiftService,
		Validate:     mockValidator,
	}

	tests := []struct {
		name           string
		jsonBody       string
		mockSetup      func()
		expectedStatus int
		expectedBody   gin.H
	}{
		{
			name: "Success",
			jsonBody: `{
				"address": "123 Main St",
				"bankName": "Bank of Test",
				"countryIso2": "US",
				"countryName": "United States",
				"isHeadquarter": true,
				"swiftCode": "ABCDEF12XXX"
			}`,
			mockSetup: func() {
				mockSwiftService.EXPECT().AddSwift(gomock.Any(), gomock.Any(), mockSwiftRepo, mockValidator).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   gin.H{"message": "Swift code added successfully"},
		},
		{
			name:     "Error - Validation failed",
			jsonBody: `dsadasdasdasd`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"message": "Bad request"},
		},
		{
			name: "Error - Validation failed",
			jsonBody: `{
				"address": "123 Main St",
				"bankName": "Bank of Test",
				"countryIso2": "US",
				"countryName": "United States",
				"isHeadquarter": "dsadas",
				"swiftCode": "ABCDEF12XXX"
			}`,
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"message": "isHeadquarter should be bool"},
		},
		{
			name: "Error - Swift code already exists",
			jsonBody: `{
				"address": "123 Main St",
				"bankName": "Bank of Test",
				"countryIso2": "US",
				"countryName": "United States",
				"isHeadquarter": true,
				"swiftCode": "ABCDEF12XXX"
			}`,
			mockSetup: func() {
				mockSwiftService.EXPECT().AddSwift(gomock.Any(), gomock.Any(), mockSwiftRepo, mockValidator).Return(customErrors.ErrSwiftCodeAlreadyExists)
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   gin.H{"message": "Swift code already exists"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/swift", strings.NewReader(tt.jsonBody))

			controller.AddSwift(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != nil {
				var responseBody gin.H
				err := json.Unmarshal(w.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, responseBody)
			}
		})
	}
}

func TestController_DeleteSwift(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSwiftRepo := mocks.NewMockSwiftRepo(ctrl)
	mockSwiftService := mocks.NewMockSwiftService(ctrl)

	controller := Controller{
		SwiftRepo:    mockSwiftRepo,
		SwiftService: mockSwiftService,
	}

	tests := []struct {
		name           string
		swiftCode      string
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:      "Success",
			swiftCode: "ABCDEF12XXX",
			mockSetup: func() {
				mockSwiftService.EXPECT().DeleteSwift(gomock.Any(), "ABCDEF12XXX", mockSwiftRepo).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "Error - Swift not found",
			swiftCode: "INVALIDCODE",
			mockSetup: func() {
				mockSwiftService.EXPECT().DeleteSwift(gomock.Any(), "INVALIDCODE", mockSwiftRepo).Return(customErrors.ErrSwiftNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodDelete, "/swift/"+tt.swiftCode, nil)
			c.Params = gin.Params{gin.Param{Key: "swiftCode", Value: tt.swiftCode}}

			controller.DeleteSwift(c)
			assert.Equal(t, tt.expectedStatus, c.Writer.Status())
		})
	}
}
