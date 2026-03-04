package validators

import (
	"fmt"
	"strings"
)

type Float32Validator struct {
	values []string
}

func NewFloat32Validator(v []any) Float32Validator {
	parts := make([]string, len(v))
	for i, v := range v {
		parts[i] = fmt.Sprintf("%f", v)
	}
	return Float32Validator{values: parts}
}

func (s Float32Validator) ValidatorFunc() string {
	return "ValidateFloat32"
}

func (s Float32Validator) ValidatorReq() string {
	return "Float32Request"
}

func (s Float32Validator) ValidatorResp() string {
	return "Float32Response"
}

func (s Float32Validator) ValidatorValue() string {
	return "ValueFloat32"
}

func (s Float32Validator) Values() string {
	return fmt.Sprintf("%s", strings.Join(s.values, ", "))
}

func (s Float32Validator) ValuesMatcher() string {
	if len(s.values) == 0 {
		return ""
	}
	return fmt.Sprintf("[]float32{%#v}", strings.Join(s.values, ", "))
}
