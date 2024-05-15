package ast

import (
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type IdentifierFormatterOpt func(*IdentifierFormatter)

type Formatter func(string) string

type IdentifierFormatter struct {
	packageFormatter    Formatter
	classFormatter      Formatter
	classFieldFormatter Formatter
	enumFormatter       Formatter
	enumMemberFormatter Formatter
	constFormatter      Formatter
	varFormatter        Formatter
}

func PackageFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.packageFormatter = formatter
	}
}

func ClassFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.classFormatter = formatter
	}
}

func ClassFieldFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.classFieldFormatter = formatter
	}
}

func EnumFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.enumFormatter = formatter
	}
}

func EnumMemberFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.enumMemberFormatter = formatter
	}
}

func ConstantFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.constFormatter = formatter
	}
}

func VariableFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.varFormatter = formatter
	}
}

func NewIdentifierFormatter(opts ...IdentifierFormatterOpt) *IdentifierFormatter {
	noopFormatter := func(input string) string { return input }
	formatter := &IdentifierFormatter{
		packageFormatter:    noopFormatter,
		classFormatter:      noopFormatter,
		enumFormatter:       noopFormatter,
		enumMemberFormatter: noopFormatter,
		constFormatter:      noopFormatter,
		varFormatter:        noopFormatter,
	}

	for _, opt := range opts {
		opt(formatter)
	}

	return formatter
}

func (pass *IdentifierFormatter) Process(schemas []*Schema) ([]*Schema, error) {
	newSchemas := Schemas(schemas).DeepCopy()

	for i, schema := range newSchemas {
		newSchemas[i] = pass.processSchema(schemas, schema)
	}

	return newSchemas, nil
}

func (pass *IdentifierFormatter) processSchema(originalSchemas Schemas, schema *Schema) *Schema {
	originalObjects := schema.Objects

	schema.Package = pass.packageFormatter(schema.Package)
	schema.EntryPoint = pass.classFormatter(schema.EntryPoint)
	schema.Objects = orderedmap.New[string, Object]()

	originalObjects.Iterate(func(_ string, object Object) {
		schema.AddObject(pass.processObject(originalSchemas, object))
	})

	return schema
}

func (pass *IdentifierFormatter) processObject(originalSchemas Schemas, object Object) Object {
	if object.Type.IsEnum() {
		object.Name = pass.enumMemberFormatter(object.Name)
	} else if object.Type.IsStruct() {
		object.Name = pass.classFormatter(object.Name)
	} else if object.Type.IsConcreteScalar() {
		object.Name = pass.constFormatter(object.Name)
	}

	object.SelfRef.ReferredPkg = pass.packageFormatter(object.SelfRef.ReferredPkg)
	object.SelfRef.ReferredType = object.Name

	object.Type = pass.processType(originalSchemas, object.Type)

	return object
}

func (pass *IdentifierFormatter) processType(originalSchemas Schemas, def Type) Type {
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

func (pass *IdentifierFormatter) processArray(originalSchemas Schemas, def Type) Type {
	def.Array.ValueType = pass.processType(originalSchemas, def.Array.ValueType)

	return def
}

func (pass *IdentifierFormatter) processMap(originalSchemas Schemas, def Type) Type {
	def.Map.IndexType = pass.processType(originalSchemas, def.Map.IndexType)
	def.Map.ValueType = pass.processType(originalSchemas, def.Map.ValueType)

	return def
}

func (pass *IdentifierFormatter) processEnum(def Type) Type {
	def.Enum.Values = tools.Map(def.Enum.Values, func(member EnumValue) EnumValue {
		member.Name = pass.enumMemberFormatter(member.Name)
		return member
	})

	return def
}

func (pass *IdentifierFormatter) processDisjunction(originalSchemas Schemas, def Type) Type {
	for i, branch := range def.Disjunction.Branches {
		def.Disjunction.Branches[i] = pass.processType(originalSchemas, branch)
	}

	return def
}

func (pass *IdentifierFormatter) processIntersection(originalSchemas Schemas, def Type) Type {
	for i, branch := range def.Intersection.Branches {
		def.Intersection.Branches[i] = pass.processType(originalSchemas, branch)
	}

	return def
}

func (pass *IdentifierFormatter) processRef(originalSchemas Schemas, def Type) Type {
	referredObj, found := originalSchemas.LocateObject(def.Ref.ReferredPkg, def.Ref.ReferredType)
	if !found {
		return def
	}

	def.Ref.ReferredPkg = pass.packageFormatter(def.Ref.ReferredPkg)

	if referredObj.Type.IsEnum() {
		def.Ref.ReferredType = pass.enumMemberFormatter(referredObj.Name)
	} else if referredObj.Type.IsStruct() {
		def.Ref.ReferredType = pass.classFormatter(referredObj.Name)
	} else if referredObj.Type.IsConcreteScalar() {
		def.Ref.ReferredType = pass.constFormatter(referredObj.Name)
	}

	return def
}

func (pass *IdentifierFormatter) processStruct(originalSchemas Schemas, def Type) Type {
	for i, field := range def.Struct.Fields {
		def.Struct.Fields[i].Name = pass.classFieldFormatter(field.Name)
		def.Struct.Fields[i].Type = pass.processType(originalSchemas, field.Type)
	}

	return def
}
