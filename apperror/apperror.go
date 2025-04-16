package apperror

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/michaelyusak/go-helper/appconstant"
)

type AppError struct {
	Code            int
	Message         string
	ResponseMessage string
}

type AppErrorOpt struct {
	Code            int
	Message         string
	ResponseMessage string
}

func (ae AppError) Error() string {
	return fmt.Sprintf("Code %d | Error: %s", ae.Code, ae.Message)
}

func NewAppError(code int, message string, responseMessage string) *AppError {
	return &AppError{
		Code:            code,
		Message:         message,
		ResponseMessage: responseMessage,
	}
}

func NotFoundError() *AppError {
	return NewAppError(http.StatusNotFound, appconstant.MsgNotFound, appconstant.MsgNotFound)
}

func BadRequestError(opt AppErrorOpt) *AppError {
	if opt.Message == "" {
		opt.Message = fmt.Sprintf("Bad Request Error | Stack: %s", string(debug.Stack()))
	}

	if opt.ResponseMessage == "" {
		opt.ResponseMessage = "bad request"
	}

	if opt.Code == 0 {
		opt.Code = http.StatusBadRequest
	}

	return NewAppError(opt.Code, opt.Message, opt.ResponseMessage)
}

func InternalServerError(opt AppErrorOpt) *AppError {
	if opt.Message == "" {
		opt.Message = fmt.Sprintf("Internal Server Error | Stack: %s", string(debug.Stack()))
	}

	if opt.ResponseMessage == "" {
		opt.ResponseMessage = "internal server error"
	}

	if opt.Code == 0 {
		opt.Code = http.StatusInternalServerError
	}

	return NewAppError(opt.Code, opt.Message, opt.ResponseMessage)
}
