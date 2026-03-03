package validators

import (
	"fmt"
	"strings"
)

type StringValidator struct {
	values []string
}

func NewStringValidator(v []any) StringValidator {
	parts := make([]string, len(v))
	for i, s := range v {
		parts[i] = fmt.Sprintf("%v", s)
	}
	return StringValidator{
		values: parts,
	}
}

func (s StringValidator) ValidatorFunc() string {
	return "ValidateString"
}

func (s StringValidator) ValidatorReq() string {
	return "StringRequest"
}

func (s StringValidator) ValidatorResp() string {
	return "StringResponse"
}

func (s StringValidator) ValidatorValue() string {
	return "ValueString"
}

func (s StringValidator) Values() string {
	return strings.Join(s.values, ", ")
}

func (s StringValidator) ValuesMatcher() string {
	return fmt.Sprintf("[]string{%#v}", strings.Join(s.values, ", "))
}
