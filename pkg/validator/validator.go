package validator

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func New() (*CustomValidator, error) {
	validator := validator.New()
	return &CustomValidator{Validator: validator}, nil
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.Validator.Struct(i); err != nil {
		validatorErrors := err.(validator.ValidationErrors)
		return customError(validatorErrors)
	}

	return nil
}

func customError(e validator.ValidationErrors) error {
	mapError := make(map[string]string)
	for _, err := range e {
		mapError[err.Field()] = errorMessage(err)
	}

	return echo.NewHTTPError(http.StatusBadRequest, mapError)
}

// Handle error message for specific validation tag
func errorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return "Invalid email format"
	case "gender":
		return "Gender must be male or female"
	case "role":
		return "Invalid role"
	case "gt":
		return fmt.Sprintf("Param %s should be greater than %v", err.Field(), err.Param())
	case "required_without":
		return fmt.Sprintf("One of %s or %s is required", err.Field(), err.Param())
	}
	return err.Error()
}
