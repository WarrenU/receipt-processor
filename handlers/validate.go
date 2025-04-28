package handlers

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func ValidateStruct(s interface{}) []*FieldError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errs []*FieldError
	for _, e := range err.(validator.ValidationErrors) {
		var message string
		switch e.Tag() {
		case "required":
			message = "is required"
		case "datetime":
			if e.Field() == "PurchaseDate" {
				message = "must match format YYYY-MM-DD"
			} else if e.Field() == "PurchaseTime" {
				message = "must match format HH:MM"
			} else {
				message = "must match datetime format"
			}
		default:
			message = "is invalid"
		}

		errs = append(errs, &FieldError{
			Field: e.Field(),
			Error: message,
		})
	}
	return errs
}
