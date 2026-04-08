package nullable_fields

import (
	"encoding/json"
	cog "github.com/grafana/cog/generated/cog"
	"errors"
	"fmt"
)

type Struct struct {
    A *MyObject `json:"a"`
    B *MyObject `json:"b,omitempty"`
    C *string `json:"c"`
    D []string `json:"d"`
    E map[string]*string `json:"e"`
    F *NullableFieldsStructF `json:"f"`
    G *string `json:"g"`
}

// NewStruct creates a new Struct object.
func NewStruct() *Struct {
	return &Struct{
		A: NewMyObject(),
		D: []string{},
		E: map[string]*string{},
		F: NewNullableFieldsStructF(),
		G: (func (input string) *string { return &input })(ConstantRef),
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `Struct` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *Struct) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "a"
	if fields["a"] != nil {
		if string(fields["a"]) != "null" {
			
			resource.A = &MyObject{}
			if err := resource.A.UnmarshalJSONStrict(fields["a"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("a", err)...)
			}
		
		}
		delete(fields, "a")
	} else {errs = append(errs, cog.MakeBuildErrors("a", errors.New("required field is missing from input"))...)
	}
	// Field "b"
	if fields["b"] != nil {
		if string(fields["b"]) != "null" {
			
			resource.B = &MyObject{}
			if err := resource.B.UnmarshalJSONStrict(fields["b"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("b", err)...)
			}
		
		}
		delete(fields, "b")
	
	}
	// Field "c"
	if fields["c"] != nil {
		if string(fields["c"]) != "null" {
			if err := json.Unmarshal(fields["c"], &resource.C); err != nil {
				errs = append(errs, cog.MakeBuildErrors("c", err)...)
			}
		
		}
		delete(fields, "c")
	} else {errs = append(errs, cog.MakeBuildErrors("c", errors.New("required field is missing from input"))...)
	}
	// Field "d"
	if fields["d"] != nil {
		if string(fields["d"]) != "null" {
			
			if err := json.Unmarshal(fields["d"], &resource.D); err != nil {
				errs = append(errs, cog.MakeBuildErrors("d", err)...)
			}
		
		}
		delete(fields, "d")
	} else {errs = append(errs, cog.MakeBuildErrors("d", errors.New("required field is missing from input"))...)
	}
	// Field "e"
	if fields["e"] != nil {
		if string(fields["e"]) != "null" {
			
			if err := json.Unmarshal(fields["e"], &resource.E); err != nil {
				errs = append(errs, cog.MakeBuildErrors("e", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("e", errors.New("required field is null"))...)
		
		}
		delete(fields, "e")
	} else {errs = append(errs, cog.MakeBuildErrors("e", errors.New("required field is missing from input"))...)
	}
	// Field "f"
	if fields["f"] != nil {
		if string(fields["f"]) != "null" {
			
			resource.F = &NullableFieldsStructF{}
			if err := resource.F.UnmarshalJSONStrict(fields["f"]); err != nil {
				errs = append(errs, cog.MakeBuildErrors("f", err)...)
			}
		
		}
		delete(fields, "f")
	} else {errs = append(errs, cog.MakeBuildErrors("f", errors.New("required field is missing from input"))...)
	}
	// Field "g"
	if fields["g"] != nil {
		if string(fields["g"]) != "null" {
			if err := json.Unmarshal(fields["g"], &resource.G); err != nil {
				errs = append(errs, cog.MakeBuildErrors("g", err)...)
			}
		
		}
		delete(fields, "g")
	} else {errs = append(errs, cog.MakeBuildErrors("g", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("Struct", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `Struct` objects.
func (resource Struct) Equals(other Struct) bool {
		if resource.A == nil && other.A != nil || resource.A != nil && other.A == nil {
			return false
		}

		if resource.A != nil {
		if !resource.A.Equals(*other.A) {
			return false
		}
		}
		if resource.B == nil && other.B != nil || resource.B != nil && other.B == nil {
			return false
		}

		if resource.B != nil {
		if !resource.B.Equals(*other.B) {
			return false
		}
		}
		if resource.C == nil && other.C != nil || resource.C != nil && other.C == nil {
			return false
		}

		if resource.C != nil {
		if *resource.C != *other.C {
			return false
		}
		}

		if len(resource.D) != len(other.D) {
			return false
		}

		for i1 := range resource.D {
		if resource.D[i1] != other.D[i1] {
			return false
		}
		}

		if len(resource.E) != len(other.E) {
			return false
		}

		for key1 := range resource.E {
		if resource.E[key1] == nil && other.E[key1] != nil || resource.E[key1] != nil && other.E[key1] == nil {
			return false
		}

		if resource.E[key1] != nil {
		if *resource.E[key1] != *other.E[key1] {
			return false
		}
		}
		}
		if resource.F == nil && other.F != nil || resource.F != nil && other.F == nil {
			return false
		}

		if resource.F != nil {
		if !resource.F.Equals(*other.F) {
			return false
		}
		}
		if resource.G == nil && other.G != nil || resource.G != nil && other.G == nil {
			return false
		}

		if resource.G != nil {
        if *resource.G != *other.G {
            return false
        }
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `Struct` fields for violations and returns them.
func (resource Struct) Validate() error {
	var errs cog.BuildErrors
		if resource.A != nil {
		if err := resource.A.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("a", err)...)
		}
		}
		if resource.B != nil {
		if err := resource.B.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("b", err)...)
		}
		}
		if resource.F != nil {
		if err := resource.F.Validate(); err != nil {
			errs = append(errs, cog.MakeBuildErrors("f", err)...)
		}
		}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


type MyObject struct {
    Field string `json:"field"`
}

// NewMyObject creates a new MyObject object.
func NewMyObject() *MyObject {
	return &MyObject{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `MyObject` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *MyObject) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "field"
	if fields["field"] != nil {
		if string(fields["field"]) != "null" {
			if err := json.Unmarshal(fields["field"], &resource.Field); err != nil {
				errs = append(errs, cog.MakeBuildErrors("field", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("field", errors.New("required field is null"))...)
		
		}
		delete(fields, "field")
	} else {errs = append(errs, cog.MakeBuildErrors("field", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("MyObject", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `MyObject` objects.
func (resource MyObject) Equals(other MyObject) bool {
		if resource.Field != other.Field {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `MyObject` fields for violations and returns them.
func (resource MyObject) Validate() error {
	return nil
}


const ConstantRef = "hey"

type NullableFieldsStructF struct {
    A string `json:"a"`
}

// NewNullableFieldsStructF creates a new NullableFieldsStructF object.
func NewNullableFieldsStructF() *NullableFieldsStructF {
	return &NullableFieldsStructF{
}
}
// UnmarshalJSONStrict implements a custom JSON unmarshalling logic to decode `NullableFieldsStructF` from JSON.
// Note: the unmarshalling done by this function is strict. It will fail over required fields being absent from the input, fields having an incorrect type, unexpected fields being present, …
func (resource *NullableFieldsStructF) UnmarshalJSONStrict(raw []byte) error {
	if raw == nil {
		return nil
	}
	var errs cog.BuildErrors

	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	// Field "a"
	if fields["a"] != nil {
		if string(fields["a"]) != "null" {
			if err := json.Unmarshal(fields["a"], &resource.A); err != nil {
				errs = append(errs, cog.MakeBuildErrors("a", err)...)
			}
		} else {errs = append(errs, cog.MakeBuildErrors("a", errors.New("required field is null"))...)
		
		}
		delete(fields, "a")
	} else {errs = append(errs, cog.MakeBuildErrors("a", errors.New("required field is missing from input"))...)
	}

	for field := range fields {
		errs = append(errs, cog.MakeBuildErrors("NullableFieldsStructF", fmt.Errorf("unexpected field '%s'", field))...)
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}


// Equals tests the equality of two `NullableFieldsStructF` objects.
func (resource NullableFieldsStructF) Equals(other NullableFieldsStructF) bool {
		if resource.A != other.A {
			return false
		}

	return true
}


// Validate checks all the validation constraints that may be defined on `NullableFieldsStructF` fields for violations and returns them.
func (resource NullableFieldsStructF) Validate() error {
	return nil
}


