package controllers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func PrintValidationErrors(err error) {
	for _, fieldErr := range err.(validator.ValidationErrors) {
		fmt.Printf("Field: %s, Value: %v, Tag: %s, Param: %s",
			fieldErr.Field(),
			fieldErr.Value(),
			fieldErr.Tag(),
			fieldErr.Param(),
		)
	}

}
