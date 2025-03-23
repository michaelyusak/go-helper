package apperror

import (
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/michaelyusak/go-helper/appconstant"
)

type AppError struct {
	Code    int
	Err     error
	Message string
	stack   []byte
}

func (ae AppError) Error() string {
	return fmt.Sprintf("Error %d %s from: %s", ae.Code, ae.Message, ae.Err.Error())
}

func (ae *AppError) GetStackTrace() []byte {
	return ae.stack
}

func NewAppError(code int, err error, message string) *AppError {
	return &AppError{
		Code:    code,
		Err:     err,
		Message: message,
		stack:   debug.Stack(),
	}
}

func NotFoundError() *AppError {
	err := errors.New(appconstant.MsgNotFound)

	return NewAppError(http.StatusNotFound, err, appconstant.MsgNotFound)
}

func BadRequestError(err error) *AppError {
	return NewAppError(http.StatusBadRequest, err, err.Error())
}

func InternalServerError(err error) *AppError {
	return NewAppError(http.StatusInternalServerError, err, appconstant.MsgInternalServerError)
}
