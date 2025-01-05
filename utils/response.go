package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// JSONResponse is a helper function to standardize API responses
func JSONResponse(c *gin.Context, status int, message string, data interface{}, errors interface{}) {
	c.JSON(status, gin.H{
		"message": message,
		"data":    data,
		"errors":  errors,
	})
}

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errors[field] = fmt.Sprintf("%s is required", field)
			case "max":
				errors[field] = fmt.Sprintf("%s exceeds maximum length of %s", field, fieldError.Param())
			case "min":
				errors[field] = fmt.Sprintf("%s is below minimum length of %s", field, fieldError.Param())
			case "gt":
				errors[field] = fmt.Sprintf("%s must be greater than %s", field, fieldError.Param())
			case "gte":
				errors[field] = fmt.Sprintf("%s must be greater than or equal to %s", field, fieldError.Param())
			case "lt":
				errors[field] = fmt.Sprintf("%s must be less than %s", field, fieldError.Param())
			case "lte":
				errors[field] = fmt.Sprintf("%s must be less than or equal to %s", field, fieldError.Param())
			case "isnumeric":
				errors[field] = fmt.Sprintf("%s must be a numeric value", field)
			default:
				errors[field] = fmt.Sprintf("Invalid value for %s", field)
			}
		}
	}

	return errors
}
