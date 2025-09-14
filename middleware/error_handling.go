package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/michaelyusak/go-helper/apperror"
	"github.com/michaelyusak/go-helper/dto"
)

func ErrorHandlerMiddleware(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		firstError := c.Errors[0].Err
		if firstError != nil {
			statusCode, errorResponse := checkError(firstError)
			c.AbortWithStatusJSON(statusCode, errorResponse)
		}
	}
}

func checkError(err error) (int, dto.ErrorResponse) {
	var ve validator.ValidationErrors

	var appErr *apperror.AppError

	var unmarshalErr json.UnmarshalTypeError
	unmarchalErrType := &unmarshalErr

	var jse json.SyntaxError
	jseErrType := &jse

	if errors.As(err, &ve) {
		details := generateValidationErrors(ve)
		return http.StatusBadRequest, dto.ErrorResponse{Message: "validation error", StatusCode: http.StatusBadRequest, Details: details}

	} else if errors.As(err, &appErr) {
		return appErr.Code, dto.ErrorResponse{Message: appErr.ResponseMessage, StatusCode: appErr.Code}

	} else if errors.As(err, &unmarchalErrType) {
		return http.StatusBadRequest, dto.ErrorResponse{Message: "unmarshal error", StatusCode: http.StatusBadRequest}

	} else if errors.As(err, &jseErrType) {
		return http.StatusBadRequest, dto.ErrorResponse{Message: "json syntax error", StatusCode: http.StatusBadRequest}
	}

	return http.StatusInternalServerError, dto.ErrorResponse{Message: "internal error", StatusCode: http.StatusInternalServerError}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "lte":
		return "should be less than or equal to " + fe.Param()
	case "gte":
		return "should be greater than or equal to " + fe.Param()
	case "max":
		return "should be max " + fe.Param() + " characters"
	case "email":
		return "should be in valid email format"
	case "ltefield":
		return "should be less than field " + fe.Param()
	case "ValidPassword":
		return "should be minimum 8 characters, contains at least 1 uppercase letter, and contains 1 number"
	case "latitude":
		return "should be in latitude format"
	case "longitude":
		return "should be in longitude format"
	case "number":
		return "should be a number"
	case "datetime":
		return "should be in valid date format yyyy-mm-dd"
	}
	return "unknown error"
}

func generateValidationErrors(ve validator.ValidationErrors) []dto.ValidationErrorDetails {
	details := make([]dto.ValidationErrorDetails, len(ve))
	for i, fe := range ve {
		details[i] = dto.ValidationErrorDetails{Field: fe.Field(), Message: getErrorMsg(fe)}
	}
	return details
}
