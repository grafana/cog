package defaults

type SomeStruct struct {
	FieldBool bool `json:"fieldBool"`
	FieldString string `json:"fieldString"`
	FieldStringWithConstantValue string `json:"FieldStringWithConstantValue"`
	FieldFloat32 float32 `json:"FieldFloat32"`
	FieldInt32 int32 `json:"FieldInt32"`
}

func (resource SomeStruct) Equals(other SomeStruct) bool {
		if resource.FieldBool != other.FieldBool {
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
		if resource.FieldInt32 != other.FieldInt32 {
			return false
		}

	return true
}


