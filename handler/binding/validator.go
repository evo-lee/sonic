package binding

import "github.com/go-playground/validator/v10"

var defaultValidator = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()
	v.SetTagName("binding")
	return v
}

func ValidatorEngine() *validator.Validate {
	return defaultValidator
}

func ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}
	return defaultValidator.Struct(obj)
}
