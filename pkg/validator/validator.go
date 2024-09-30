package validator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	validator_lib "github.com/go-playground/validator/v10"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
)

const (
	INVALID_NAME_CHARS = "*/!@#$%^&*()_+-={}[];:'\"?><.,|\\"
)

type validate struct {
	validator *validator_lib.Validate
}

func NewValidator() datatypes.Validator {
	vald := validator_lib.New()

	vald.RegisterValidation("name", func(fl validator_lib.FieldLevel) bool { return true })

	for i := 0; i < 1000; i += 5 {

		vald.RegisterValidation("size:"+strconv.Itoa(i), func(fl validator_lib.FieldLevel) bool { return true })
	}
	return &validate{
		validator: vald,
	}
}

func (v *validate) Struct(st interface{}) (fieldName string, err error) {
	val := reflect.ValueOf(st)
	if val.Kind() != reflect.Struct {
		log.Fatalln("type of st is not struct")
	}

	for i := 0; i < val.NumField(); i++ {
		// fieldName := val.Type().Field(i).Name
		fieldValue := fmt.Sprintf("%v", val.Field(i))
		fieldTag := val.Type().Field(i).Tag.Get("validate")
		fieldNameJson := val.Type().Field(i).Tag.Get("json")

		if err := v.ValidateField(fieldValue, fieldTag); err != nil {
			return fieldNameJson, err
		}

	}

	return "", nil
}

func (v *validate) ValidateField(fieldValue any, tag string) error {
	err := v.validator.Var(fieldValue, tag)

	if strings.Contains(tag, "omitempty") {
		if fieldValue == "" {
			return nil
		}

	} else if strings.Contains(tag, "e164") {
		if err != nil {
			return errors.New("invalid phone number")
		}
	} else if strings.Contains(tag, "name") {

		for char := range INVALID_NAME_CHARS {
			if strings.Contains(fieldValue.(string), string(char)) {
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
