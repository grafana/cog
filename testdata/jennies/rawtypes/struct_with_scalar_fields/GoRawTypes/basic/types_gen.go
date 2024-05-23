package basic

// This
// is
// a
// comment
type SomeStruct struct {
	// Anything can go in there.
// Really, anything.
FieldAny any `json:"FieldAny"`
	FieldBool bool `json:"FieldBool"`
	FieldBytes []byte `json:"FieldBytes"`
	FieldString string `json:"FieldString"`
	FieldStringWithConstantValue string `json:"FieldStringWithConstantValue"`
	FieldFloat32 float32 `json:"FieldFloat32"`
	FieldFloat64 float64 `json:"FieldFloat64"`
	FieldUint8 uint8 `json:"FieldUint8"`
	FieldUint16 uint16 `json:"FieldUint16"`
	FieldUint32 uint32 `json:"FieldUint32"`
	FieldUint64 uint64 `json:"FieldUint64"`
	FieldInt8 int8 `json:"FieldInt8"`
	FieldInt16 int16 `json:"FieldInt16"`
	FieldInt32 int32 `json:"FieldInt32"`
	FieldInt64 int64 `json:"FieldInt64"`
}

