package controllers

import (
	"awesomeProject/customErrors"
	"awesomeProject/repositories"
	"awesomeProject/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	BankRepo repositories.BankRepo
}

func (controller Controller) GetDetails(c *gin.Context) {
	ctx := c.Request.Context()
	swiftCode := c.Param("swiftCode")
	bank, branches, err := services.GetBankDetails(ctx, swiftCode, controller.BankRepo)
	if err != nil {
		var bankErr *customErrors.HttpError
		if errors.As(err, &bankErr) {
			bankErr.Send(c)
			return
		}
		customErrors.ErrUnknown.Send(c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"address":       bank.Address,
		"bankName":      bank.Name,
		"countryISO2":   bank.CountryIso2,
		"countryName":   bank.CountryName,
		"isHeadquarter": bank.IsHeadquarter,
		"swiftCode":     bank.SwiftCode,
		"branches":      branches,
	})
}
