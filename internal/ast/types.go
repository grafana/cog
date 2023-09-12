package ast

type Kind string

const (
	KindDisjunction Kind = "disjunction"
	KindRef         Kind = "ref"

	KindStruct Kind = "struct"
	KindEnum   Kind = "enum"
	KindMap    Kind = "map"

	KindArray Kind = "array"

	KindScalar Kind = "scalar"
)

type ScalarKind string

const (
	KindNull   ScalarKind = "null"
	KindAny    ScalarKind = "any"
	KindBytes  ScalarKind = "bytes"
	KindString ScalarKind = "string"

	KindFloat32 ScalarKind = "float32"
	KindFloat64 ScalarKind = "float64"

	KindUint8  ScalarKind = "uint8"
	KindUint16 ScalarKind = "uint16"
	KindUint32 ScalarKind = "uint32"
	KindUint64 ScalarKind = "uint64"
	KindInt8   ScalarKind = "int8"
	KindInt16  ScalarKind = "int16"
	KindInt32  ScalarKind = "int32"
	KindInt64  ScalarKind = "int64"

	KindBool ScalarKind = "bool"
)

type TypeConstraint struct {
	// TODO: something more descriptive here? constant?
	Op   string
	Args []any
}

// Struct representing every type defined by the IR.
// Bonus: in a way that can be (un)marshaled to/from JSON,
// which is useful for unit tests.
type Type struct {
	Kind Kind

	Disjunction *DisjunctionType `json:",omitempty"`
	Array       *ArrayType       `json:",omitempty"`
	Enum        *EnumType        `json:",omitempty"`
	Map         *MapType         `json:",omitempty"`
	Struct      *StructType      `json:",omitempty"`
	Ref         *RefType         `json:",omitempty"`
	Scalar      *ScalarType      `json:",omitempty"`
}

func Any() Type {
	return NewScalar(KindAny)
}

func Null() Type {
	return NewScalar(KindNull)
}

func Bool() Type {
	return NewScalar(KindBool)
}

func Bytes() Type {
	return NewScalar(KindBytes)
}

func String() Type {
	return NewScalar(KindString)
}

func NewDisjunction(branches Types) Type {
	return Type{
		Kind: KindDisjunction,
		Disjunction: &DisjunctionType{
			Branches: branches,
		},
	}
}

func NewArray(valueType Type) Type {
	return Type{
		Kind: KindArray,
		Array: &ArrayType{
			ValueType: valueType,
		},
	}
}

func NewEnum(values []EnumValue) Type {
	return Type{
		Kind: KindEnum,
		Enum: &EnumType{
			Values: values,
		},
	}
}

func NewMap(indexType Type, valueType Type) Type {
	return Type{
		Kind: KindMap,
		Map: &MapType{
			IndexType: indexType,
			ValueType: valueType,
		},
	}
}

func NewStruct(fields []StructField) Type {
	return Type{
		Kind: KindStruct,
		Struct: &StructType{
			Fields: fields,
		},
	}
}

func NewRef(referredTypeName string) Type {
	return Type{
		Kind: KindRef,
		Ref: &RefType{
			ReferredType: referredTypeName,
		},
	}
}

func NewScalar(kind ScalarKind) Type {
	return Type{
		Kind: KindScalar,
		Scalar: &ScalarType{
			ScalarKind: kind,
		},
	}
}

func (t Type) IsNull() bool {
	return t.Kind == KindScalar && t.AsScalar().ScalarKind == KindNull
}

func (t Type) IsAny() bool {
	return t.Kind == KindScalar && t.AsScalar().ScalarKind == KindAny
}

func (t Type) AsDisjunction() DisjunctionType {
	return *t.Disjunction
}

func (t Type) AsArray() ArrayType {
	return *t.Array
}

func (t Type) AsEnum() EnumType {
	return *t.Enum
}

func (t Type) AsMap() MapType {
	return *t.Map
}

func (t Type) AsStruct() StructType {
	return *t.Struct
}

func (t Type) AsRef() RefType {
	return *t.Ref
}

func (t Type) AsScalar() ScalarType {
	return *t.Scalar
}

// named declaration of a type
type Object struct {
	Name     string
	Comments []string
	Type     Type
}

type File struct { //nolint: musttag
	Package     string
	Definitions []Object
}

func (file *File) LocateDefinition(name string) Object {
	for _, def := range file.Definitions {
		if def.Name == name {
			return def
		}
	}

	return Object{}
}

type Types []Type

func (types Types) HasNullType() bool {
	for _, t := range types {
		if t.IsNull() {
			return true
		}
	}

	return false
}

func (types Types) NonNullTypes() Types {
	results := make([]Type, 0, len(types))

	for _, t := range types {
		if t.IsNull() {
			continue
		}

		results = append(results, t)
	}

	return results
}

type DisjunctionType struct {
	Branches Types
}

type ArrayType struct {
	ValueType Type
}

type EnumType struct {
	Values []EnumValue // possible values. Value types might be different
}

type EnumValue struct {
	Type  Type
	Name  string
	Value any
}

type MapType struct {
	IndexType Type
	ValueType Type
}

type StructType struct {
	Fields []StructField
}

type StructField struct {
	Name     string
	Comments []string
	Type     Type
	Required bool
	Default  any
}

type RefType struct {
	ReferredType string
}

type ScalarType struct {
	ScalarKind  ScalarKind // bool, bytes, string, int, float, ...
	Value       any        // if value isn't nil, we're representing a constant scalar
	Constraints []TypeConstraint
}
