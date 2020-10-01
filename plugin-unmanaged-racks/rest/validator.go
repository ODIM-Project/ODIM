package rest

import (
	"fmt"
	"strings"
)

type ValidationResult struct {
	errors map[string]string
}

func (ve *ValidationResult) RegisterError(fieldName, errorMessage string) {
	ve.errors[fieldName] = errorMessage
}

func (ve *ValidationResult) Error() string {
	var messages []string
	for k, v := range ve.errors {
		messages = append(messages, fmt.Sprintf("%s: %v", k, v))
	}

	return strings.Join(messages, "#")
}

func (ve *ValidationResult) HasErrors() bool {
	return len(ve.errors) > 0
}

type compositeValidator []validator

func (v *compositeValidator) validate() ValidationResult {
	vr := ValidationResult{
		errors: make(map[string]string),
	}
	for _, validationStep := range *v {
		if validationStep.validationRule() {
			vr.RegisterError(validationStep.field, validationStep.message)
		}
	}
	return vr
}

type validationRule func() bool

type validator struct {
	validationRule
	field, message string
}
