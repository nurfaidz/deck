package helpers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"strings"
)

func TranslateErrorMessage(err error) map[string]string {
	errorsMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errorsMap[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errorsMap[field] = fmt.Sprintf("%s is not a valid email", field)
			case "unique":
				errorsMap[field] = fmt.Sprintf("%s is already exists", field)
			case "min":
				errorsMap[field] = fmt.Sprintf("%s must be at least %s character", field, fieldError.Param())
			case "max":
				errorsMap[field] = fmt.Sprintf("%s must be at most %s character", field, fieldError.Param())
			case "numeric":
				errorsMap[field] = fmt.Sprintf("%s must be numeric", field)
			case "oneof":
				options := strings.Split(fieldError.Param(), ",")
				errorsMap[field] = fmt.Sprintf("%s must be one of the following: %s", field, strings.Join(options, ", "))
			default:
				errorsMap[field] = "Invalid value"

			}
		}
	}

	if err != nil {
		// file is not found
		if strings.Contains(err.Error(), "http: no such file") || strings.Contains(err.Error(), "missing file") {
			errorsMap["image"] = "Image file is required"
		}

		if strings.Contains(err.Error(), "Duplicate entry") {
			if strings.Contains(err.Error(), "username") {
				errorsMap["username"] = "Username already exists"
			}
			if strings.Contains(err.Error(), "email") {
				errorsMap["email"] = "Email already exists"
			}
		} else if err == gorm.ErrRecordNotFound {
			errorsMap["record"] = "Record not found"
		}
	}

	return errorsMap
}

func IsDuplicateEntryError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Duplicate entry")
}
