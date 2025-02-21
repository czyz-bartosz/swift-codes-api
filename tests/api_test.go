package tests

import (
	"awesomeProject/configs"
	"awesomeProject/controllers"
	"awesomeProject/dbs"
	"awesomeProject/dbs/migrations"
	"awesomeProject/models"
	"awesomeProject/repositories"
	"awesomeProject/routes"
	"awesomeProject/services"
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestEnvironment(t *testing.T) (*bun.DB, *controllers.Controller) {
	err := godotenv.Load("../.env")

	if err != nil {
		t.Fatalf("Env file error: %v", err)
	}

	config := configs.GetConfig()
	config.DBConfig.Host = "localhost"

	db := dbs.Connect(&config.DBConfig)

	_, err = db.NewDropTable().IfExists().Model((*models.Swift)(nil)).Exec(context.Background())
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}

	err = migrations.Migrate(db)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	swiftRepo := &repositories.SwiftRepoPostgres{
		Db: &dbs.BunDBWrapper{DB: db},
	}

	validate := validator.New()
	validate.RegisterStructValidation(models.SwiftStructLevelValidation, models.Swift{})

	swiftService := services.SwiftServiceDefault{}

	swiftController := controllers.Controller{
		SwiftService: &swiftService,
		SwiftRepo:    swiftRepo,
		Validate:     validate,
	}

	return db, &swiftController
}

func TestAddSwift(t *testing.T) {
	db, swiftController := setupTestEnvironment(t)
	defer db.Close()

	router := routes.SetupRouter(swiftController)

	server := httptest.NewServer(router)
	defer server.Close()

	swift := models.Swift{
		SwiftCode:     "XXXXXXXXXXX",
		BankName:      "Test Bank",
		Address:       "123 Test Street",
		CountryIso2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}

	jsonData, err := json.Marshal(swift)
	assert.NoError(t, err)

	resp, err := http.Post(server.URL+"/v1/swift-codes/", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "Swift code added successfully", response["message"])
}

func TestGetSwiftDetails(t *testing.T) {
	db, swiftController := setupTestEnvironment(t)
	defer db.Close()

	router := routes.SetupRouter(swiftController)

	server := httptest.NewServer(router)
	defer server.Close()

	swift := models.Swift{
		SwiftCode:     "XXXXXXXXXXX",
		BankName:      "Test Bank",
		Address:       "123 Test Street",
		CountryIso2:   "US",
		CountryName:   "UNITED STATES",
		IsHeadquarter: true,
	}

	jsonData, err := json.Marshal(swift)
	assert.NoError(t, err)

	resp, err := http.Post(server.URL+"/v1/swift-codes/", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(server.URL + "/v1/swift-codes/" + swift.SwiftCode)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, swift.SwiftCode, response["swiftCode"])
	assert.Equal(t, swift.BankName, response["bankName"])
	assert.Equal(t, swift.Address, response["address"])
	assert.Equal(t, swift.CountryIso2, response["countryISO2"])
	assert.Equal(t, swift.CountryName, response["countryName"])
	assert.Equal(t, swift.IsHeadquarter, response["isHeadquarter"])
}

func TestDeleteSwift(t *testing.T) {
	db, swiftController := setupTestEnvironment(t)
	defer db.Close()

	router := routes.SetupRouter(swiftController)

	server := httptest.NewServer(router)
	defer server.Close()

	swift := models.Swift{
		SwiftCode:     "XXXXXXXXXXX",
		BankName:      "Test Bank",
		Address:       "123 Test Street",
		CountryIso2:   "US",
		CountryName:   "United States",
		IsHeadquarter: true,
	}

	jsonData, err := json.Marshal(swift)
	assert.NoError(t, err)

	resp, err := http.Post(server.URL+"/v1/swift-codes/", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	req, err := http.NewRequest("DELETE", server.URL+"/v1/swift-codes/"+swift.SwiftCode, nil)
	assert.NoError(t, err)

	client := &http.Client{}
	resp, err = client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	resp, err = http.Get(server.URL + "/v1/swift-codes/" + swift.SwiftCode)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetSwiftsDetailsByCountryIso2Code(t *testing.T) {
	db, swiftController := setupTestEnvironment(t)
	defer db.Close()

	router := routes.SetupRouter(swiftController)

	server := httptest.NewServer(router)
	defer server.Close()

	swift := models.Swift{
		SwiftCode:     "XXXXXXXXXXX",
		BankName:      "Test Bank",
		Address:       "123 Test Street",
		CountryIso2:   "US",
		CountryName:   "UNITED STATES",
		IsHeadquarter: true,
	}

	jsonData, err := json.Marshal(swift)
	assert.NoError(t, err)

	resp, err := http.Post(server.URL+"/v1/swift-codes/", "application/json", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(server.URL + "/v1/swift-codes/country/US")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, "US", response["countryISO2"])
	assert.Equal(t, "UNITED STATES", response["countryName"])
	assert.NotEmpty(t, response["swiftCodes"])
}
