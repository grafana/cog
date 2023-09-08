package ast

type Kind string

const (
	KindDisjunction Kind = "disjunction"
	KindRef         Kind = "ref"

	KindStruct Kind = "struct"
	KindEnum   Kind = "enum"
	KindMap    Kind = "map"

	KindNull   Kind = "null"
	KindAny    Kind = "any"
	KindBytes  Kind = "bytes"
	KindArray  Kind = "array"
	KindString Kind = "string"

	KindFloat32 Kind = "float32"
	KindFloat64 Kind = "float64"

	KindUint8  Kind = "uint8"
	KindUint16 Kind = "uint16"
	KindUint32 Kind = "uint32"
	KindUint64 Kind = "uint64"
	KindInt8   Kind = "int8"
	KindInt16  Kind = "int16"
	KindInt32  Kind = "int32"
	KindInt64  Kind = "int64"

	KindBool Kind = "bool"
)

type TypeConstraint struct {
	// TODO: something more descriptive here? constant?
	Op   string
	Args []any
}

// interface for every type that we can represent:
// struct, enum, array, string, int, ...
type Type interface {
	Kind() Kind
}

// named declaration of a type
type Object struct {
	Name     string
	Comments []string
	Type     Type
}

type File struct {
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

var _ Type = (*DisjunctionType)(nil)

type Types []Type

func (types Types) HasNullType() bool {
	for _, t := range types {
		if t.Kind() == KindNull {
			return true
		}
	}

	return false
}

func (types Types) NonNullTypes() Types {
	results := make([]Type, 0, len(types))

	for _, t := range types {
		if t.Kind() == KindNull {
			continue
		}

		results = append(results, t)
	}

	return results
}

type DisjunctionType struct {
	Branches Types
}

func (disjunctionType DisjunctionType) Kind() Kind {
	return KindDisjunction
}

var _ Type = (*ArrayType)(nil)

type ArrayType struct {
	ValueType Type
}

func (arrayType ArrayType) Kind() Kind {
	return KindArray
}

var _ Type = (*EnumType)(nil)

type EnumType struct {
	Values []EnumValue // possible values. Value types might be different
}

type EnumValue struct {
	Type  Type
	Name  string
	Value any
}

func (arrayType EnumType) Kind() Kind {
	return KindEnum
}

var _ Type = (*MapType)(nil)

type MapType struct {
	IndexType Type
	ValueType Type
}

func (arrayType MapType) Kind() Kind {
	return KindMap
}

var _ Type = (*StructType)(nil)

type StructType struct {
	Fields []StructField
}

func (structType StructType) Kind() Kind {
	return KindStruct
}

type StructField struct {
	Name     string
	Comments []string
	Type     Type
	Required bool
	Default  any
}

var _ Type = (*RefType)(nil)

type RefType struct {
	ReferredType string
}

func (refType RefType) Kind() Kind {
	return KindRef
}

var _ Type = (*ScalarType)(nil)

type ScalarType struct {
	ScalarKind  Kind // bool, bytes, string, int, float, ...
	Value       any  // if value isn't nil, we're representing a constant scalar
	Constraints []TypeConstraint
}

func (scalarType ScalarType) Kind() Kind {
	return scalarType.ScalarKind
}
