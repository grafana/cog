package basic

// This
// is
// a
// comment
type SomeStruct struct {
	// Anything can go in there.
// Really, anything.
FieldAny any `json:"FieldAny" yaml:"FieldAny"`
	FieldBool bool `json:"FieldBool" yaml:"FieldBool"`
	FieldBytes []byte `json:"FieldBytes" yaml:"FieldBytes"`
	FieldString string `json:"FieldString" yaml:"FieldString"`
	FieldStringWithConstantValue string `json:"FieldStringWithConstantValue" yaml:"FieldStringWithConstantValue"`
	FieldFloat32 float32 `json:"FieldFloat32" yaml:"FieldFloat32"`
	FieldFloat64 float64 `json:"FieldFloat64" yaml:"FieldFloat64"`
	FieldUint8 uint8 `json:"FieldUint8" yaml:"FieldUint8"`
	FieldUint16 uint16 `json:"FieldUint16" yaml:"FieldUint16"`
	FieldUint32 uint32 `json:"FieldUint32" yaml:"FieldUint32"`
	FieldUint64 uint64 `json:"FieldUint64" yaml:"FieldUint64"`
	FieldInt8 int8 `json:"FieldInt8" yaml:"FieldInt8"`
	FieldInt16 int16 `json:"FieldInt16" yaml:"FieldInt16"`
	FieldInt32 int32 `json:"FieldInt32" yaml:"FieldInt32"`
	FieldInt64 int64 `json:"FieldInt64" yaml:"FieldInt64"`
}

func (resource SomeStruct) Equals(other SomeStruct) bool {
		// is DeepEqual good enough here?
		if !reflect.DeepEqual(resource.FieldAny, other.FieldAny) {
			return false
		}
		if resource.FieldBool != other.FieldBool {
			return false
		}
		if resource.FieldBytes != other.FieldBytes {
			return false
		}
		if resource.FieldString != other.FieldString {
			return false
		}
		if resource.FieldStringWithConstantValue != other.FieldStringWithConstantValue {
			return false
		}
		if resource.FieldFloat32 != other.FieldFloat32 {
			return false
		}
		if resource.FieldFloat64 != other.FieldFloat64 {
			return false
		}
		if resource.FieldUint8 != other.FieldUint8 {
			return false
		}
		if resource.FieldUint16 != other.FieldUint16 {
			return false
		}
		if resource.FieldUint32 != other.FieldUint32 {
			return false
		}
		if resource.FieldUint64 != other.FieldUint64 {
			return false
		}
		if resource.FieldInt8 != other.FieldInt8 {
			return false
		}
		if resource.FieldInt16 != other.FieldInt16 {
			return false
		}
		if resource.FieldInt32 != other.FieldInt32 {
			return false
		}
		if resource.FieldInt64 != other.FieldInt64 {
			return false
		}

	return true
}


