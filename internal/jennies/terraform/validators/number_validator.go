package validators

import (
	"fmt"
	"strings"
)

// NumberValidator is used for int8, uint8, int16, uint16
type NumberValidator struct {
	values []string
}

func NewNumberValidator(v []any) NumberValidator {
	parts := make([]string, len(v))
	for i, s := range v {
		parts[i] = fmt.Sprintf("%f", s)
	}
	return NumberValidator{values: parts}
}

func (s NumberValidator) ValidatorFunc() string {
	return "ValidateNumber"
}

func (s NumberValidator) ValidatorReq() string {
	return "NumberRequest"
}

func (s NumberValidator) ValidatorResp() string {
	return "NumberResponse"
}

func (s NumberValidator) ValidatorValue() string {
	return "ValueBigFloat"
}

func (s NumberValidator) Values() string {
	return strings.Join(s.values, ", ")
}

func (s NumberValidator) ValuesMatcher() string {
	return fmt.Sprintf("[]int{%#v}", strings.Join(s.values, ", "))
}
