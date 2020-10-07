package redfish

type ValidationResult []MsgExtendedInfo

func (ve *ValidationResult) Error() *CommonError {
	e := NewError()
	for _, info := range *ve {
		e.AddExtendedInfo(info)
	}
	return e
}

func (ve *ValidationResult) HasErrors() bool {
	return len(*ve) > 0
}

type CompositeValidator []Validator

func (v *CompositeValidator) Validate() ValidationResult {
	vr := ValidationResult{}
	for _, validationStep := range *v {
		if validationStep.ValidationRule() {
			vr = append(vr, validationStep.ErrorGenerator())
		}
	}
	return vr
}

type ValidationRule func() bool
type ErrorGenerator func() MsgExtendedInfo

type Validator struct {
	ValidationRule
	ErrorGenerator
}
