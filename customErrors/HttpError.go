package customErrors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpError struct {
	code    int
	message string
}

func (e *HttpError) Error() string {
	return e.message
}

func (e *HttpError) Code() int { return e.code }

func (e *HttpError) Message() string { return e.message }

func (e *HttpError) Send(c *gin.Context) {
	c.JSON(e.Code(), gin.H{"message": e.Message()})
}

func NewHttpError(code int, message string) *HttpError {
	return &HttpError{
		code:    code,
		message: message,
	}
}

var ErrBankNotFound = NewHttpError(http.StatusNotFound, "Bank not found")
var ErrUnknown = NewHttpError(http.StatusInternalServerError, "Something went wrong")
