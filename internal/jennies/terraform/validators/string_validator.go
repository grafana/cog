package validators

import (
	"fmt"
	"strings"
)

type StringValidator struct {
	values []string
}

func NewStringValidator(v []any) StringValidator {
	if len(v) == 0 {
		return StringValidator{}
	}
	parts := make([]string, 0, len(v))
	for _, s := range v {
		if s == nil {
			continue
		}
		parts = append(parts, fmt.Sprintf("%v", s))
	}

	if len(parts) == 0 {
		return StringValidator{}
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
	if len(s.values) == 0 {
		return ""
	}
	return fmt.Sprintf("%#v", s.values)
}
