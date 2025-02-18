package controllers

import (
	"awesomeProject/customErrors"
	"awesomeProject/models"
	"awesomeProject/repositories"
	"awesomeProject/services"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	BankRepo repositories.BankRepo
}

func handleError(c *gin.Context, err error) {
	var bankErr *customErrors.HttpError
	if errors.As(err, &bankErr) {
		bankErr.Send(c)
		return
	}
	customErrors.ErrUnknown.Send(c)
}

func (controller Controller) GetBankDetails(c *gin.Context) {
	ctx := c.Request.Context()
	swiftCode := c.Param("swiftCode")
	bank, branches, err := services.GetBankDetails(ctx, swiftCode, controller.BankRepo)
	if err != nil {
		handleError(c, err)
		return
	}

	if models.IsSwiftCodeOfHeadquarter(swiftCode) {
		c.JSON(http.StatusOK, gin.H{
			"address":       bank.Address,
			"bankName":      bank.Name,
			"countryISO2":   bank.CountryIso2,
			"countryName":   bank.CountryName,
			"isHeadquarter": bank.IsHeadquarter,
			"swiftCode":     bank.SwiftCode,
			"branches":      branches,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"address":       bank.Address,
			"bankName":      bank.Name,
			"countryISO2":   bank.CountryIso2,
			"countryName":   bank.CountryName,
			"isHeadquarter": bank.IsHeadquarter,
			"swiftCode":     bank.SwiftCode,
		})
	}
}

func (controller Controller) GetBanksDetailsByCountryIso2Code(c *gin.Context) {
	ctx := c.Request.Context()
	countryIso2Code := c.Param("countryIso2Code")

	countryName, banks, err := services.GetBanksDetailsByCountryIso2Code(ctx, countryIso2Code, controller.BankRepo)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"countryISO2": countryIso2Code,
		"countryName": countryName,
		"swiftCodes":  banks,
	})
}

func (controller Controller) AddBank(c *gin.Context) {
	ctx := c.Request.Context()
	var bank models.Bank
	err := c.ShouldBindJSON(&bank)

	if err != nil {
		var typeError *json.UnmarshalTypeError
		if errors.As(err, &typeError) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": typeError.Field + " should be " + typeError.Type.Name(),
			})
			return
		}
		handleError(c, customErrors.ErrBadRequest)
		return
	}

	err = services.AddBank(ctx, &bank, controller.BankRepo)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Swift code added successfully",
	})
}
