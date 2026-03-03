package validators

import (
	"fmt"
	"strings"
)

type BoolValidator struct {
	values []string
}

func NewBoolValidator(v []any) BoolValidator {
	parts := make([]string, len(v))
	for i, v := range v {
		parts[i] = fmt.Sprintf("%t", v)
	}
	return BoolValidator{values: parts}
}

func (s BoolValidator) ValidatorFunc() string {
	return "ValidateBool"
}

func (s BoolValidator) ValidatorReq() string {
	return "BoolRequest"
}

func (s BoolValidator) ValidatorResp() string {
	return "BoolResponse"
}

func (s BoolValidator) ValidatorValue() string {
	return "ValueBool"
}

func (s BoolValidator) Values() string {
	return strings.Join(s.values, ", ")
}

func (s BoolValidator) ValuesMatcher() string {
	return fmt.Sprintf("[]bool{%#v}", strings.Join(s.values, ", "))
}
