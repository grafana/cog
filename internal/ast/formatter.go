package ast

import (
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type IdentifierFormatterOpt func(*IdentifierFormatter)

type Formatter func(string) string

type IdentifierFormatter struct {
	Package     Formatter
	Object      Formatter
	ObjectField Formatter
	Enum        Formatter
	EnumMember  Formatter
	Constant    Formatter
	Variable    Formatter
}

func NewIdentifierFormatter(opts ...IdentifierFormatterOpt) *IdentifierFormatter {
	noopFormatter := func(input string) string { return input }
	formatter := &IdentifierFormatter{
		Package:    noopFormatter,
		Object:     noopFormatter,
		Enum:       noopFormatter,
		EnumMember: noopFormatter,
		Constant:   noopFormatter,
		Variable:   noopFormatter,
	}

	for _, opt := range opts {
		opt(formatter)
	}

	return formatter
}

func PackageFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.Package = formatter
	}
}

func ObjectFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.Object = formatter
	}
}

func ObjectFieldFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.ObjectField = formatter
	}
}

func EnumFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.Enum = formatter
	}
}

func EnumMemberFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.EnumMember = formatter
	}
}

func ConstantFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.Constant = formatter
	}
}

func VariableFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.Variable = formatter
	}
}

type IdentifierFormatterPass struct {
	formatter *IdentifierFormatter
}

func NewIdentifierFormatterPass(formatter *IdentifierFormatter) *IdentifierFormatterPass {
	return &IdentifierFormatterPass{
		formatter: formatter,
	}
}

func (pass *IdentifierFormatterPass) Process(schemas []*Schema) ([]*Schema, error) {
	newSchemas := Schemas(schemas).DeepCopy()

	for i, schema := range newSchemas {
		newSchemas[i] = pass.processSchema(schemas, schema)
	}

	return newSchemas, nil
}

func (pass *IdentifierFormatterPass) processSchema(originalSchemas Schemas, schema *Schema) *Schema {
	originalObjects := schema.Objects

	schema.Package = pass.formatter.Package(schema.Package)
	schema.EntryPoint = pass.formatter.Object(schema.EntryPoint)
	schema.Objects = orderedmap.New[string, Object]()

	originalObjects.Iterate(func(_ string, object Object) {
		schema.AddObject(pass.processObject(originalSchemas, object))
	})

	return schema
}

func (pass *IdentifierFormatterPass) processObject(originalSchemas Schemas, object Object) Object {
	if object.Type.IsEnum() {
		object.Name = pass.formatter.EnumMember(object.Name)
	} else if object.Type.IsStruct() {
		object.Name = pass.formatter.Object(object.Name)
	} else if object.Type.IsConcreteScalar() {
		object.Name = pass.formatter.Constant(object.Name)
	}

	object.SelfRef.ReferredPkg = pass.formatter.Package(object.SelfRef.ReferredPkg)
	object.SelfRef.ReferredType = object.Name

	object.Type = pass.processType(originalSchemas, object.Type)

	return object
}

func (pass *IdentifierFormatterPass) processType(originalSchemas Schemas, def Type) Type {
	if def.IsArray() {
		return pass.processArray(originalSchemas, def)
	}

	if def.IsMap() {
		return pass.processMap(originalSchemas, def)
	}

	if def.IsEnum() {
		return pass.processEnum(def)
	}

	if def.IsStruct() {
		return pass.processStruct(originalSchemas, def)
	}

	if def.IsDisjunction() {
		return pass.processDisjunction(originalSchemas, def)
	}

	if def.IsIntersection() {
		return pass.processIntersection(originalSchemas, def)
	}

	if def.IsRef() {
		return pass.processRef(originalSchemas, def)
	}

	return def
}

func (pass *IdentifierFormatterPass) processArray(originalSchemas Schemas, def Type) Type {
	def.Array.ValueType = pass.processType(originalSchemas, def.Array.ValueType)

	return def
}

func (pass *IdentifierFormatterPass) processMap(originalSchemas Schemas, def Type) Type {
	def.Map.IndexType = pass.processType(originalSchemas, def.Map.IndexType)
	def.Map.ValueType = pass.processType(originalSchemas, def.Map.ValueType)

	return def
}

func (pass *IdentifierFormatterPass) processEnum(def Type) Type {
	def.Enum.Values = tools.Map(def.Enum.Values, func(member EnumValue) EnumValue {
		member.Name = pass.formatter.EnumMember(member.Name)
		return member
	})

	return def
}

func (pass *IdentifierFormatterPass) processDisjunction(originalSchemas Schemas, def Type) Type {
	for i, branch := range def.Disjunction.Branches {
		def.Disjunction.Branches[i] = pass.processType(originalSchemas, branch)
	}

	return def
}

func (pass *IdentifierFormatterPass) processIntersection(originalSchemas Schemas, def Type) Type {
	for i, branch := range def.Intersection.Branches {
		def.Intersection.Branches[i] = pass.processType(originalSchemas, branch)
	}

	return def
}

func (pass *IdentifierFormatterPass) processRef(originalSchemas Schemas, def Type) Type {
	referredObj, found := originalSchemas.LocateObject(def.Ref.ReferredPkg, def.Ref.ReferredType)
	if !found {
		return def
	}

	def.Ref.ReferredPkg = pass.formatter.Package(def.Ref.ReferredPkg)

	if referredObj.Type.IsEnum() {
		def.Ref.ReferredType = pass.formatter.Enum(referredObj.Name)
	} else if referredObj.Type.IsStruct() {
		def.Ref.ReferredType = pass.formatter.Object(referredObj.Name)
	} else if referredObj.Type.IsConcreteScalar() {
		def.Ref.ReferredType = pass.formatter.Constant(referredObj.Name)
	}

	return def
}

func (pass *IdentifierFormatterPass) processStruct(originalSchemas Schemas, def Type) Type {
	for i, field := range def.Struct.Fields {
		def.Struct.Fields[i].Name = pass.formatter.ObjectField(field.Name)
		def.Struct.Fields[i].Type = pass.processType(originalSchemas, field.Type)
	}

	return def
}
