package validator

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"reflect"
	"regexp"
	"strings"

	validator_lib "github.com/go-playground/validator/v10"
)

var (
	INVALID_NAME_CHARS        = []string{"*", "/", "!", "@", "#", "$%", "^", "&", "*", "(", ")", "_", "+", "-", "=", "{", "}", "[", "]", ";", ":", "'", "\\", "?", "\"", ".", ">", "<", "|", "union"}
	INVALID_DESCRIPTION_CHARS = []string{"<", ">", "'"}
	INVALID_USERAGENT_CHARS   = []string{"\\", "union", "<", "'"}
)

type validate struct {
	validator *validator_lib.Validate
}

func NewValidator() *validate {
	vald := validator_lib.New()

	vald.RegisterValidation("name", func(fl validator_lib.FieldLevel) bool {
		for _, char := range INVALID_NAME_CHARS {
			if strings.Contains(fmt.Sprintf("%v", fl.Field()), char) {
				return false
			}
		}
		return true
	})
	vald.RegisterValidation("description", func(fl validator_lib.FieldLevel) bool {
		for _, char := range INVALID_DESCRIPTION_CHARS {
			if strings.Contains(fmt.Sprintf("%v", fl.Field()), char) {
				return false
			}
		}
		return true
	})
	vald.RegisterValidation("phone_number", func(fl validator_lib.FieldLevel) bool {
		ok, err := regexp.Match("\\+?9\\d{9}$", []byte(fl.Field().String()))
		if err != nil {
			slog.Error("error in number regex", "err", err)
		}

		if len(fl.Field().String()) != 13 {
			return false
		}

		return ok
	})
	vald.RegisterValidation("username", func(fl validator_lib.FieldLevel) bool {
		ok, err := regexp.Match("^([A-Za-z]|\\.){1,200}$", []byte(fl.Field().String()))
		if err != nil {
			slog.Error("error in number regex", "err", err)
		}

		if len(fl.Field().String()) != 13 {
			return false
		}

		return ok
	})
	vald.RegisterValidation("useragent", func(fl validator_lib.FieldLevel) bool {
		for _, char := range INVALID_USERAGENT_CHARS {
			if strings.Contains(fmt.Sprintf("%v", fl.Field()), char) {
				return false
			}
		}

		return true
	})

	return &validate{
		validator: vald,
	}
}

func (v *validate) Struct(st interface{}) (errMap map[string]string) {
	val := reflect.ValueOf(st)
	if val.Kind() != reflect.Struct {
		log.Fatalln("type of st is not struct")
	}

	errMap = make(map[string]string)

	errs := v.validator.Struct(st)
	if errs != nil {
		for _, err := range errs.(validator_lib.ValidationErrors) {
			errMap[err.Field()] = err.Error()
		}
	}

	if len(errMap) == 0 {
		return nil

	} else {
		return errMap
	}
}

func (v *validate) ValidateFieldByFieldName(fieldName string, fieldValue any, model any) error {
	sf, ok := reflect.TypeOf(model).FieldByName(fieldName)
	if !ok {
		panic("cannot find fieldByName:" + fieldName)
	}

	validateTag, ok := sf.Tag.Lookup("validate")
	if !ok {
		// no validation setted.
		return nil
	}

	return v.ValidateField(fieldValue, validateTag)
}

func (v *validate) ValidateField(fieldValue any, tag string) error {
	err := v.validator.Var(fieldValue, tag)
	// fmt.Println("----------", fieldValue, tag, fieldValue == "")
	// fmt.Println("value", fieldValue, "tag", tag)
	if strings.Contains(tag, "allowempty") {
		if fieldValue == "" {
			return nil
		}

		// } else if strings.Contains(tag, "e164") {
		// 	if err != nil {
		// 		return errors.New("invalid phone number")
		// 	}
		// } else if strings.Contains(tag, "name") {

		// 	for char := range INVALID_NAME_CHARS {
		// 		if strings.Contains(fieldValue.(string), string(char)) {
		// 			return errors.New("invalid character: " + string(char))
		// 		}
		// 	}
	} else if strings.Contains(tag, "required") {
		if fieldValue == "" {
			return errors.New("this field is requried")
		}

	}

	if err != nil {
		return errors.New("invalid field")
	} else {
		return nil
	}

}
