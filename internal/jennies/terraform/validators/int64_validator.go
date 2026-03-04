package validators

import (
	"fmt"
	"strings"
)

type Int64Validator struct {
	values []string
}

func NewInt64Validator(v []any) Int64Validator {
	parts := make([]string, len(v))
	for i, s := range v {
		parts[i] = fmt.Sprintf("%d", s)
	}
	return Int64Validator{values: parts}
}

func (s Int64Validator) ValidatorFunc() string {
	return "ValidateInt64"
}

func (s Int64Validator) ValidatorReq() string {
	return "Int64Request"
}

func (s Int64Validator) ValidatorResp() string {
	return "Int64Response"
}

func (s Int64Validator) ValidatorValue() string {
	return "ValueInt64"
}

func (s Int64Validator) Values() string {
	if len(s.values) == 0 {
		return ""
	}
	return strings.Join(s.values, ", ")
}

func (s Int64Validator) ValuesMatcher() string {
	return fmt.Sprintf("[]int64{%#v}", strings.Join(s.values, ", "))
}
