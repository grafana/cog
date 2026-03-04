package validators

import (
	"fmt"
	"strings"
)

type Int32Validator struct {
	values []string
}

func NewInt32Validator(v []any) Int32Validator {
	parts := make([]string, len(v))
	for i, s := range v {
		parts[i] = fmt.Sprintf("%d", s)
	}
	return Int32Validator{values: parts}
}

func (s Int32Validator) ValidatorFunc() string {
	return "ValidateInt32"
}

func (s Int32Validator) ValidatorReq() string {
	return "Int32Request"
}

func (s Int32Validator) ValidatorResp() string {
	return "Int32Response"
}

func (s Int32Validator) ValidatorValue() string {
	return "ValueInt32"
}

func (s Int32Validator) Values() string {
	return strings.Join(s.values, ", ")
}

func (s Int32Validator) ValuesMatcher() string {
	if len(s.values) == 0 {
		return ""
	}
	return fmt.Sprintf("[]int32{%#v}", strings.Join(s.values, ", "))
}
