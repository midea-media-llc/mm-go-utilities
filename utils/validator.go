package utils

import (
	"futa.express.api.accountant/utils/logs"
	"github.com/go-playground/validator/v10"
)

// ValidationError is implementation of validation error
type ValidationError struct {
	MessageCode string
	Data        map[string]string
}

func (r ValidationError) Error() string {
	return r.MessageCode
}

// Validate date using defined tag
func Validate(data interface{}) *ValidationError {
	var validate *validator.Validate
	validate = validator.New()

	err := validate.Struct(data)

	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			logs.Errorf(err.Error())
			return nil
		}

		for _, err := range err.(validator.ValidationErrors) {

			logs.Infof(err.Namespace())
			logs.Infof(err.Field())
			logs.Infof(err.StructNamespace())
			logs.Infof(err.StructField())
			logs.Infof(err.Tag())
			logs.Infof(err.ActualTag())
			logs.Infof("%v", err.Kind())
			logs.Infof("%v", err.Type())
			logs.Infof("%v", err.Value())
			logs.Infof(err.Param())
			logs.Infof("%v", err)

			return &ValidationError{
				MessageCode: getMessageCodeByValidateTag(err.Tag()),
				Data: map[string]string{
					"Name":  err.StructField(),
					"Param": err.Param(),
				},
			}
		}
	}

	return nil
}

func getMessageCodeByValidateTag(tag string) string {
	switch tag {
	case "required":
		return "validation_required"
	case "max":
		return "validation_max"
	case "min":
		return "validation_min"
	case "oneof":
		return "validation_oneof"
	}

	return "validation_default"
}
