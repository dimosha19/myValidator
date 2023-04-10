package validator

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var IsDigit = regexp.MustCompile(`^[-]*[0-9]+$`).MatchString

type IValidator interface {
	Validate() ValidationErrors
}

type IntValidator struct {
	name string
	val  any
	tag  string
}

type StrValidator struct {
	name string
	val  any
	tag  string
}

func makeerr(name string, e error) ValidationError {
	if e.Error() == ErrInvalidValidatorSyntax.Error() {
		return ValidationError{e}
	}
	s := name + ": " + e.Error()
	return ValidationError{errors.New(s)}
}

func (v StrValidator) Validate() ValidationErrors {

	value := v.val.(string)
	errs := ValidationErrors{}
	elems := strings.Split(v.tag, " ")
	for i := range elems {
		pair := strings.Split(elems[i], ":")
		switch pair[0] {
		case "max":
			if err := v.ValidateMax(value, pair[1]); err != nil {
				errs = append(errs, makeerr(v.name, err))
				return errs
			}
		case "min":
			if err := v.ValidateMin(value, pair[1]); err != nil {
				errs = append(errs, makeerr(v.name, err))
				return errs
			}
		case "in":
			if err := v.ValidateIn(value, pair[1]); err != nil {
				errs = append(errs, makeerr(v.name, err))
				return errs
			}
		case "len":
			if err := v.ValidateLen(value, pair[1]); err != nil {
				errs = append(errs, makeerr(v.name, err))
				return errs
			}
		default:
			errs = append(errs, makeerr(v.name, ErrInvalidValidatorSyntax))
			return errs
		}
	}
	return errs
}

func (v StrValidator) ValidateIn(val string, tagVal string) error {
	if tagVal == "" {
		return ErrInvalidValidatorSyntax
	}
	tagArr := strings.Split(tagVal, ",")
	for i := range tagArr {
		if strings.Contains(tagArr[i], val) {
			return nil
		}
	}
	return errors.New("field does not contain required value")
}

func (v StrValidator) ValidateMax(val string, tagVal string) error {
	if !IsDigit(tagVal) {
		return ErrInvalidValidatorSyntax
	}
	elem, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if len(val) <= elem {
		return nil
	}
	return errors.New("field does not fit according to the restriction from above")
}

func (v StrValidator) ValidateMin(val string, tagVal string) error {
	elem, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if len(val) >= elem {
		return nil
	}
	return errors.New("the field does not fit the limit below")
}

func (v StrValidator) ValidateLen(val string, tagVal string) error {
	elem, err := strconv.Atoi(tagVal)
	if err != nil || elem < 0 {
		return ErrInvalidValidatorSyntax
	}
	if len(val) == elem {
		return nil
	}
	return errors.New("field has an invalid length")
}

func (v IntValidator) ValidateIn(val int, tagVal string) error {
	if tagVal == "" {
		return ErrInvalidValidatorSyntax
	}
	tagArr := strings.Split(tagVal, ",")
	for i := range tagArr {
		elem, err := strconv.Atoi(tagArr[i])
		if err != nil {
			return ErrInvalidValidatorSyntax
		}
		if val == elem {
			return nil
		}
	}
	return errors.New("field does not contain required value")
}

func (v IntValidator) ValidateMin(val int, tagVal string) error {
	elem, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if val >= elem {
		return nil
	}
	return errors.New("the field does not fit the limit below")
}

func (v IntValidator) ValidateMax(val int, tagVal string) error {
	if !IsDigit(tagVal) {
		return ErrInvalidValidatorSyntax
	}
	elem, err := strconv.Atoi(tagVal)
	if err != nil {
		return ErrInvalidValidatorSyntax
	}
	if val <= elem {
		return nil
	}
	return errors.New("field does not fit according to the restriction from above")
}

func (v IntValidator) Validate() ValidationErrors {
	value := v.val.(int)
	errs := ValidationErrors{}
	elems := strings.Split(v.tag, " ")
	for i := range elems {
		pair := strings.Split(elems[i], ":")
		switch pair[0] {
		case "max":
			if err := v.ValidateMax(value, pair[1]); err != nil {
				errs = append(errs, makeerr(v.name, err))
				return errs
			}
		case "min":
			if err := v.ValidateMin(value, pair[1]); err != nil {
				errs = append(errs, makeerr(v.name, err))
				return errs
			}
		case "in":
			if err := v.ValidateIn(value, pair[1]); err != nil {
				errs = append(errs, makeerr(v.name, err))
				return errs
			}
		default:
			errs = append(errs, makeerr(v.name, ErrInvalidValidatorSyntax))
			return errs
		}
	}
	return errs
}
