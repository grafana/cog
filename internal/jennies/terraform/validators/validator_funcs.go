package validators

type ValidatorFuncs interface {
	ValidatorFunc() string
	ValidatorReq() string
	ValidatorResp() string
	ValidatorValue() string
	Values() string
	ValuesMatcher() string
}
