package validators

import (
	"fmt"
	"strings"
)

type Float64Validator struct {
	values []string
}

func NewFloat64Validator(v []any) Float64Validator {
	parts := make([]string, len(v))
	for i, v := range v {
		parts[i] = fmt.Sprintf("%v", v)
	}
	return Float64Validator{values: parts}
}

func (s Float64Validator) ValidatorFunc() string {
	return "ValidateFloat64"
}

func (s Float64Validator) ValidatorReq() string {
	return "Float64Request"
}

func (s Float64Validator) ValidatorResp() string {
	return "Float64Response"
}

func (s Float64Validator) ValidatorValue() string {
	return "ValueFloat64"
}

func (s Float64Validator) Values() string {
	return strings.Join(s.values, ", ")
}

func (s Float64Validator) ValuesMatcher() string {
	if len(s.values) == 0 {
		return ""
	}
	return fmt.Sprintf("[]float64{%#v}", strings.Join(s.values, ", "))
}
