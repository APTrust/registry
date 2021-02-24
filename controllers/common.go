package controllers

import (
	"github.com/APTrust/registry/models"
	"github.com/go-playground/validator/v10"
)

// ValidationErrors returns a map of validation error messages
// suitable for display on an HTML form.
func ValidationErrors(err error, obj models.Model) map[string]string {
	if _, ok := err.(validator.ValidationErrors); ok {
		errMap := make(map[string]interface{})
		// First, make a map of which fields had bad values, and
		// what the bad values were. In this map, the key is the
		// name of the struct property, which usually begins with
		// an upper-case letter, not the name of the form input,
		// which usually begins with a lower-case. The value is the
		// actual (invalid) value that the user submitted. E.g.:
		//
		// {
		//    "Name":  "",
		//    "Email": "user@",
		// }
		for _, fieldErr := range err.(validator.ValidationErrors) {
			errMap[fieldErr.Field()] = fieldErr.Value()
		}
		// The model's GetValidationErrors method converts vague
		// validation errors to user-friendly messages that we
		// can display on the HTML form.
		return obj.GetValidationErrors(errMap)
	}
	return nil
}
