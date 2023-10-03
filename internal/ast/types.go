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

type Op string

const (
	MinLengthOp        Op = "minLength"
	MaxLengthOp        Op = "maxLength"
	MultipleOfOp       Op = "multipleOf"
	EqualOp            Op = "=="
	NotEqualOp         Op = "!="
	LessThanOp         Op = "<"
	LessThanEqualOp    Op = "<="
	GreaterThanOp      Op = ">"
	GreaterThanEqualOp Op = ">="
)

type TypeConstraint struct {
	Op   Op
	Args []any
}

func (constraint TypeConstraint) DeepCopy() TypeConstraint {
	newConstraint := TypeConstraint{
		Op:   constraint.Op,
		Args: make([]any, 0, len(constraint.Args)),
	}

	newConstraint.Args = append(newConstraint.Args, constraint.Args...)

	return newConstraint
}

// Struct representing every type defined by the IR.
// Bonus: in a way that can be (un)marshaled to/from JSON,
// which is useful for unit tests.
type Type struct {
	Kind     Kind
	Nullable bool
	Default  any

	Disjunction *DisjunctionType `json:",omitempty"`
	Array       *ArrayType       `json:",omitempty"`
	Enum        *EnumType        `json:",omitempty"`
	Map         *MapType         `json:",omitempty"`
	Struct      *StructType      `json:",omitempty"`
	Ref         *RefType         `json:",omitempty"`
	Scalar      *ScalarType      `json:",omitempty"`
}

func (t Type) DeepCopy() Type {
	newType := Type{
		Kind:     t.Kind,
		Nullable: t.Nullable,
		Default:  t.Default,
	}

	if t.Disjunction != nil {
		newDisjunction := t.Disjunction.DeepCopy()
		newType.Disjunction = &newDisjunction
	}
	if t.Array != nil {
		newArray := t.Array.DeepCopy()
		newType.Array = &newArray
	}
	if t.Enum != nil {
		newEnum := t.Enum.DeepCopy()
		newType.Enum = &newEnum
	}
	if t.Map != nil {
		newMap := t.Map.DeepCopy()
		newType.Map = &newMap
	}
	if t.Struct != nil {
		newStruct := t.Struct.DeepCopy()
		newType.Struct = &newStruct
	}
	if t.Ref != nil {
		newRef := t.Ref.DeepCopy()
		newType.Ref = &newRef
	}
	if t.Scalar != nil {
		newScalar := t.Scalar.DeepCopy()
		newType.Scalar = &newScalar
	}

	return newType
}

type TypeOption func(def *Type)

func Nullable() TypeOption {
	return func(def *Type) {
		def.Nullable = true
	}
}

func Default(value any) TypeOption {
	return func(def *Type) {
		def.Default = value
	}
}

func Value(value any) TypeOption {
	return func(def *Type) {
		if def.Kind != KindScalar {
			return
		}

		def.Scalar.Value = value
	}
}

func Any() Type {
	return NewScalar(KindAny)
}

func Null() Type {
	return NewScalar(KindNull)
}

func Bool(opts ...TypeOption) Type {
	return NewScalar(KindBool, opts...)
}

func Bytes(opts ...TypeOption) Type {
	return NewScalar(KindBytes, opts...)
}

func String(opts ...TypeOption) Type {
	return NewScalar(KindString, opts...)
}

func NewDisjunction(branches Types, opts ...TypeOption) Type {
	def := Type{
		Kind: KindDisjunction,
		Disjunction: &DisjunctionType{
			Branches:             branches,
			DiscriminatorMapping: make(map[string]any),
		},
	}

	for _, opt := range opts {
		opt(&def)
	}

	return def
}

func NewArray(valueType Type, opts ...TypeOption) Type {
	def := Type{
		Kind: KindArray,
		Array: &ArrayType{
			ValueType: valueType,
		},
	}

	for _, opt := range opts {
		opt(&def)
	}

	return def
}

func NewEnum(values []EnumValue, opts ...TypeOption) Type {
	def := Type{
		Kind: KindEnum,
		Enum: &EnumType{
			Values: values,
		},
	}

	for _, opt := range opts {
		opt(&def)
	}

	return def
}

func NewMap(indexType Type, valueType Type, opts ...TypeOption) Type {
	def := Type{
		Kind: KindMap,
		Map: &MapType{
			IndexType: indexType,
			ValueType: valueType,
		},
	}

	for _, opt := range opts {
		opt(&def)
	}

	return def
}

func NewStruct(fields ...StructField) Type {
	return Type{
		Kind: KindStruct,
		Struct: &StructType{
			Fields: fields,
			Hint:   make(map[JennyHint]any),
		},
	}
}

func NewRef(referredPkg string, referredTypeName string, opts ...TypeOption) Type {
	def := Type{
		Kind: KindRef,
		Ref: &RefType{
			ReferredPkg:  referredPkg,
			ReferredType: referredTypeName,
		},
	}

	for _, opt := range opts {
		opt(&def)
	}

	return def
}

func NewScalar(kind ScalarKind, opts ...TypeOption) Type {
	def := Type{
		Kind: KindScalar,
		Scalar: &ScalarType{
			ScalarKind: kind,
		},
	}

	for _, opt := range opts {
		opt(&def)
	}

	return def
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

func NewObject(name string, objectType Type) Object {
	return Object{
		Name: name,
		Type: objectType,
	}
}

func (object Object) DeepCopy() Object {
	newObject := Object{
		Name: object.Name,
		Type: object.Type.DeepCopy(),
	}

	newObject.Comments = append(newObject.Comments, object.Comments...)

	return newObject
}

type Files []*File

func (files Files) DeepCopy() []*File {
	newFiles := make([]*File, 0, len(files))

	for _, file := range files {
		newFile := file.DeepCopy()
		newFiles = append(newFiles, &newFile)
	}

	return newFiles
}

type File struct { //nolint: musttag
	Package     string
	Definitions []Object
}

func (file *File) DeepCopy() File {
	newFile := File{
		Package:     file.Package,
		Definitions: make([]Object, 0, len(file.Definitions)),
	}

	for _, def := range file.Definitions {
		newFile.Definitions = append(newFile.Definitions, def.DeepCopy())
	}

	return newFile
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

func (types Types) HasOnlyScalarOrArray() bool {
	for _, t := range types {
		if t.Kind == KindArray {
			if !t.AsArray().IsArrayOfScalars() {
				return false
			}

			continue
		}

		if t.Kind != KindScalar {
			return false
		}
	}

	return true
}

func (types Types) HasOnlyRefs() bool {
	for _, t := range types {
		if t.Kind != KindRef {
			return false
		}
	}

	return true
}

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

	// If the branches are references to structs, some languages will need
	// extra context to be able to distinguish between them. Golang, for
	// example, doesn't support sum types (disjunctions of fixed types).
	// To emulate sum types for these languages, we need a way to
	// discriminate against every possible type.
	//
	// To do that, we need two things:
	//	- a discriminator: the name of a field that is present in all types.
	//	  The value of which identifies the type being used.
	//  - a mapping: associating a type name to its "discriminator value".
	Discriminator        string
	DiscriminatorMapping map[string]any // likely a map[string]string or map[string]int
}

func (t DisjunctionType) DeepCopy() DisjunctionType {
	newT := DisjunctionType{
		Branches:             make([]Type, 0, len(t.Branches)),
		Discriminator:        t.Discriminator,
		DiscriminatorMapping: make(map[string]any, len(t.DiscriminatorMapping)),
	}

	for _, branch := range t.Branches {
		newT.Branches = append(newT.Branches, branch.DeepCopy())
	}

	for k, v := range t.DiscriminatorMapping {
		newT.DiscriminatorMapping[k] = v
	}

	return newT
}

type ArrayType struct {
	ValueType Type
}

func (t ArrayType) DeepCopy() ArrayType {
	return ArrayType{
		ValueType: t.ValueType.DeepCopy(),
	}
}

func (t ArrayType) IsArrayOfScalars() bool {
	if t.ValueType.Kind == KindArray {
		return t.ValueType.AsArray().IsArrayOfScalars()
	}

	return t.ValueType.Kind == KindScalar
}

type EnumType struct {
	Values []EnumValue // possible values. Value types might be different
}

func (t EnumType) DeepCopy() EnumType {
	newT := EnumType{
		Values: make([]EnumValue, 0, len(t.Values)),
	}

	for _, value := range t.Values {
		newT.Values = append(newT.Values, value.DeepCopy())
	}

	return newT
}

type EnumValue struct {
	Type    Type
	Name    string
	Value   any
	Default any
}

func (t EnumValue) DeepCopy() EnumValue {
	return EnumValue{
		Type:  t.Type.DeepCopy(),
		Name:  t.Name,
		Value: t.Value,
	}
}

type MapType struct {
	IndexType Type
	ValueType Type
}

func (t MapType) DeepCopy() MapType {
	return MapType{
		IndexType: t.IndexType.DeepCopy(),
		ValueType: t.ValueType.DeepCopy(),
	}
}

type StructType struct {
	Fields []StructField
	Hint   map[JennyHint]any // Hints meant to be used by jennies
}

func (structType StructType) IsGeneratedFromDisjunction() bool {
	return structType.Hint[HintDisjunctionOfScalars] != nil ||
		structType.Hint[HintDiscriminatedDisjunctionOfRefs] != nil
}

func (structType StructType) DeepCopy() StructType {
	newT := StructType{
		Fields: make([]StructField, 0, len(structType.Fields)),
		Hint:   make(map[JennyHint]any, len(structType.Hint)),
	}

	for _, field := range structType.Fields {
		newT.Fields = append(newT.Fields, field.DeepCopy())
	}
	for k, v := range structType.Hint {
		newT.Hint[k] = v
	}

	return newT
}

func (structType StructType) FieldByName(name string) (StructField, bool) {
	// FIXME: we don't have a way to directly get a struct field by name :(
	for _, field := range structType.Fields {
		if field.Name != name {
			continue
		}

		return field, true
	}

	return StructField{}, false
}

type StructField struct {
	Name     string
	Comments []string
	Type     Type
	Required bool
}

func (structField StructField) DeepCopy() StructField {
	newT := StructField{
		Name:     structField.Name,
		Type:     structField.Type.DeepCopy(),
		Required: structField.Required,
	}

	newT.Comments = append(newT.Comments, structField.Comments...)

	return newT
}

type StructFieldOption func(field *StructField)

func Required() StructFieldOption {
	return func(field *StructField) {
		field.Required = true
	}
}

func NewStructField(name string, fieldType Type, opts ...StructFieldOption) StructField {
	field := StructField{
		Name: name,
		Type: fieldType,
	}

	for _, opt := range opts {
		opt(&field)
	}

	return field
}

type RefType struct {
	ReferredPkg  string
	ReferredType string
}

func (t RefType) DeepCopy() RefType {
	return RefType{
		ReferredPkg:  t.ReferredPkg,
		ReferredType: t.ReferredType,
	}
}

type ScalarType struct {
	ScalarKind  ScalarKind // bool, bytes, string, int, float, ...
	Value       any        // if value isn't nil, we're representing a constant scalar
	Constraints []TypeConstraint
}

func (scalarType ScalarType) DeepCopy() ScalarType {
	newT := ScalarType{
		ScalarKind:  scalarType.ScalarKind,
		Value:       scalarType.Value,
		Constraints: make([]TypeConstraint, 0, len(scalarType.Constraints)),
	}

	for _, constraint := range scalarType.Constraints {
		newT.Constraints = append(newT.Constraints, constraint.DeepCopy())
	}

	return newT
}

func (scalarType ScalarType) IsConcrete() bool {
	return scalarType.Value != nil
}
