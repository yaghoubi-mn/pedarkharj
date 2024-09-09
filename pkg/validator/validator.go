package validator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	validator_lib "github.com/go-playground/validator/v10"
)

const (
	INVALID_NAME_CHARS = "*/!@#$%^&*()_+-={}[];:'\"?><.,|\\"
)

type validate struct {
	validator *validator_lib.Validate
}

func NewValidator() validate {
	return validate{
		validator: validator_lib.New(),
	}
}

func (v *validate) Struct(st interface{}) map[string]string {
	errMap := make(map[string]string)
	val := reflect.ValueOf(st)
	if val.Kind() != reflect.Struct {
		log.Fatalln("type of st is not struct")
	}

	for i := 0; i < val.NumField(); i++ {
		fieldName := val.Type().Field(i).Name
		fieldValue := fmt.Sprintf("%v", val.Field(i))
		fieldTag := val.Type().Field(i).Tag.Get("validate")

		if err := v.ValidateField(fieldValue, fieldTag); err != nil {
			errMap[fieldName] = err.Error()
		}

	}

	return errMap
}

func (v *validate) ValidateField(fieldValue, tag string) error {
	err := v.validator.Var(fieldValue, tag)

	if strings.Contains(tag, "e164") {
		if err != nil {
			return errors.New("invalid phone number")
		}
	} else if strings.Contains(tag, "name") {

		for char := range INVALID_NAME_CHARS {
			if strings.Contains(fieldValue, string(char)) {
				return errors.New("invalid character: " + string(char))
			}
		}
	} else {
		if err != nil {
			return errors.New("invalid field")
		}
	}

	return nil
}
