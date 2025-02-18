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
	SwiftRepo repositories.SwiftRepo
}

func handleError(c *gin.Context, err error) {
	var httpErr *customErrors.HttpError
	if errors.As(err, &httpErr) {
		httpErr.Send(c)
		return
	}
	customErrors.ErrUnknown.Send(c)
}

func (controller Controller) GetSwiftDetails(c *gin.Context) {
	ctx := c.Request.Context()
	swiftCode := c.Param("swiftCode")
	swift, branches, err := services.GetSwiftDetails(ctx, swiftCode, controller.SwiftRepo)
	if err != nil {
		handleError(c, err)
		return
	}

	if models.IsSwiftCodeOfHeadquarter(swiftCode) {
		c.JSON(http.StatusOK, gin.H{
			"address":       swift.Address,
			"bankName":      swift.BankName,
			"countryISO2":   swift.CountryIso2,
			"countryName":   swift.CountryName,
			"isHeadquarter": swift.IsHeadquarter,
			"swiftCode":     swift.SwiftCode,
			"branches":      branches,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"address":       swift.Address,
			"bankName":      swift.BankName,
			"countryISO2":   swift.CountryIso2,
			"countryName":   swift.CountryName,
			"isHeadquarter": swift.IsHeadquarter,
			"swiftCode":     swift.SwiftCode,
		})
	}
}

func (controller Controller) GetSwiftsDetailsByCountryIso2Code(c *gin.Context) {
	ctx := c.Request.Context()
	countryIso2Code := c.Param("countryIso2Code")

	countryName, swifts, err := services.GetSwiftsDetailsByCountryIso2Code(ctx, countryIso2Code, controller.SwiftRepo)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"countryISO2": countryIso2Code,
		"countryName": countryName,
		"swiftCodes":  swifts,
	})
}

func (controller Controller) AddSwift(c *gin.Context) {
	ctx := c.Request.Context()
	var swift models.Swift
	err := c.ShouldBindJSON(&swift)

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

	err = services.AddSwift(ctx, &swift, controller.SwiftRepo)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Swift code added successfully",
	})
}
