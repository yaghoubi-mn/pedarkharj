package validator

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	validator_lib "github.com/go-playground/validator/v10"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
)

var (
	INVALID_NAME_CHARS      = []string{"*", "/", "!", "@", "#", "$%", "^", "&", "*", "(", ")", "_", "+", "-", "=", "{", "}", "[", "]", ";", ":", "'", "\\", "?", "\"", ".", ">", "<", "|", "union"}
	INVALID_USERAGENT_CHARS = []string{"\\", "union", "<", "'"}
)

type validate struct {
	validator *validator_lib.Validate
}

func NewValidator() datatypes.Validator {
	vald := validator_lib.New()

	vald.RegisterValidation("name", func(fl validator_lib.FieldLevel) bool {
		for _, char := range INVALID_NAME_CHARS {
			if strings.Contains(fmt.Sprintf("%v", fl.Field()), char) {
				return false
			}
		}
		fmt.Println("name", fl.Field().String(), "is valid")
		return true

	})
	vald.RegisterValidation("allowempty", func(fl validator_lib.FieldLevel) bool { return true })
	vald.RegisterValidation("phone_number", func(fl validator_lib.FieldLevel) bool {
		ok, err := regexp.Match("\\+?9\\d{9}$", []byte(fl.Field().String()))
		if err != nil {
			slog.Error("error in number regex", "err", err)
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

	for i := 0; i < 1000; i += 5 {

		vald.RegisterValidation("size:"+strconv.Itoa(i), func(fl validator_lib.FieldLevel) bool {
			return len(fl.Field().String()) <= i
		})
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
