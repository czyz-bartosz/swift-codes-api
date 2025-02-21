package controllers

import (
	"awesomeProject/customErrors"
	"awesomeProject/models"
	"awesomeProject/repositories"
	"awesomeProject/services"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Controller struct {
	SwiftRepo    repositories.SwiftRepo
	Validate     models.SwiftValidator
	SwiftService services.SwiftService
}

func handleError(c *gin.Context, err error) {
	var httpErr *customErrors.HttpError
	if errors.As(err, &httpErr) {
		httpErr.Send(c)
		return
	}
	var validationErr validator.ValidationErrors
	if errors.As(err, &validationErr) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": strings.Split(validationErr.Error(), "\n"),
		})
		return
	}
	customErrors.ErrUnknown.Send(c)
}

func (controller Controller) GetSwiftDetails(c *gin.Context) {
	ctx := c.Request.Context()
	swiftCode := c.Param("swiftCode")
	swift, branches, err := controller.SwiftService.GetSwiftDetails(ctx, swiftCode, controller.SwiftRepo)
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

	countryName, swifts, err := controller.SwiftService.GetSwiftsDetailsByCountryIso2Code(ctx, countryIso2Code, controller.SwiftRepo)
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

	err = controller.SwiftService.AddSwift(ctx, &swift, controller.SwiftRepo, controller.Validate)
	if err != nil {
		handleError(c, err)
		fmt.Print(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Swift code added successfully",
	})
}

func (controller Controller) DeleteSwift(c *gin.Context) {
	ctx := c.Request.Context()
	swiftCode := c.Param("swiftCode")
	err := controller.SwiftService.DeleteSwift(ctx, swiftCode, controller.SwiftRepo)
	if err != nil {
		handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
