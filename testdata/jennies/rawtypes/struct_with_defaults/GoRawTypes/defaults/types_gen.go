package defaults

type SomeStruct struct {
	FieldBool bool `json:"fieldBool"`
	FieldString string `json:"fieldString"`
	FieldStringWithConstantValue string `json:"FieldStringWithConstantValue"`
	FieldFloat32 float32 `json:"FieldFloat32"`
	FieldInt32 int32 `json:"FieldInt32"`
}

