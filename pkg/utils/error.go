package utils

import "github.com/go-playground/validator/v10"

type HttpError struct {
	Status int         `json:"status"`
	Field  string      `json:"field"`
	Type   string      `json:"type"`
	Kind   string      `json:"kind"`
	Value  interface{} `json:"value"`
}

func BuildError(errors validator.ValidationErrors) []HttpError {
	var errs []HttpError
	for _, err := range errors {
		error := HttpError{
			Status: 400,
			Field:  err.Field(),
			Type:   err.Type().String(),
			Kind:   err.Kind().String(),
			Value:  err.Value(),
		}
		errs = append(errs, error)
	}
	return errs
}
