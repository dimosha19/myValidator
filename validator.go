package validator

import (
	"reflect"

	"github.com/pkg/errors"
)

var ErrNotStruct = errors.New("wrong argument given, should be a struct")
var ErrInvalidValidatorSyntax = errors.New("invalid validator syntax")
var ErrValidateForUnexportedFields = errors.New("validation for unexported field is not allowed")

var tag = "validate"

type ValidationError struct {
	Err error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	res := ""
	for i := range v {
		res += v[i].Err.Error()
	}
	return res
}

func fieldProcessing(curr reflect.StructField, value reflect.Value, i int) ValidationErrors {
	errorst := ValidationErrors{}
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}

	var validator IValidator

	if curr.Type.Kind() == reflect.Int {
		validator = IntValidator{curr.Name, int(value.Field(i).Int()), curr.Tag.Get(tag)}
	} else if curr.Type.Kind() == reflect.String {
		validator = StrValidator{curr.Name, value.Field(i).String(), curr.Tag.Get(tag)}
	} else if curr.Type.Kind() == reflect.Slice && value.Field(i).Type().Elem().Kind() == reflect.Int {
		validator = IntSliceValidator{curr.Name, value.Field(i).Interface(), curr.Tag.Get(tag)}
	} else if curr.Type.Kind() == reflect.Slice && value.Field(i).Type().Elem().Kind() == reflect.String {
		validator = StrSliceValidator{curr.Name, value.Field(i).Interface(), curr.Tag.Get(tag)}
	}

	if res := validator.Validate(); res != nil {
		errorst = append(errorst, res...)
	}
	return errorst
}

func Validate(v any) error {
	errorst := ValidationErrors{}

	typ := reflect.TypeOf(v)
	value := reflect.ValueOf(v)

	if typ.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < typ.NumField(); i++ {
		curr := typ.Field(i)

		if curr.Tag.Get(tag) == "" {
			continue
		}

		if !curr.IsExported() {
			errorst = append(errorst, ValidationError{ErrValidateForUnexportedFields})
			return errorst
		}

		if a := fieldProcessing(curr, value, i); len(a) != 0 {
			errorst = append(errorst, a...)
		}
	}
	if len(errorst) != 0 {
		return errorst
	}
	return nil
}
