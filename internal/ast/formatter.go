package ast

import (
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

// region Identifier formatter

type IdentifierFormatterOpt func(*IdentifierFormatter)

type Formatter func(string) string

type IdentifierFormatter struct {
	Package            Formatter
	Object             Formatter
	ObjectPublicField  Formatter
	ObjectPrivateField Formatter
	Enum               Formatter
	EnumMember         Formatter
	Constant           Formatter
	Variable           Formatter
	Option             Formatter
}

func NewIdentifierFormatter(opts ...IdentifierFormatterOpt) *IdentifierFormatter {
	noopFormatter := func(input string) string { return input }
	formatter := &IdentifierFormatter{
		Package:            noopFormatter,
		Object:             noopFormatter,
		ObjectPublicField:  noopFormatter,
		ObjectPrivateField: noopFormatter,
		Enum:               noopFormatter,
		EnumMember:         noopFormatter,
		Constant:           noopFormatter,
		Variable:           noopFormatter,
		Option:             noopFormatter,
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

func ObjectPublicFieldFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.ObjectPublicField = formatter
	}
}

func ObjectPrivateFieldFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.ObjectPrivateField = formatter
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

func OptionFormatter(formatter Formatter) IdentifierFormatterOpt {
	return func(pass *IdentifierFormatter) {
		pass.Option = formatter
	}
}

// endregion

// region Identifier formatter compiler pass

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
		def.Struct.Fields[i] = pass.processStructField(originalSchemas, field)
	}

	return def
}

func (pass *IdentifierFormatterPass) processStructField(originalSchemas Schemas, field StructField) StructField {
	field.Name = pass.formatter.ObjectPublicField(field.Name)
	field.Type = pass.processType(originalSchemas, field.Type)

	return field
}

//endregion

// region Identifier formatter builder rewriter

type IdentifierFormatterBuilderRewriter struct {
	formatter     *IdentifierFormatter
	formatterPass *IdentifierFormatterPass
}

func NewIdentifierFormatterBuilderRewriter(formatter *IdentifierFormatter) *IdentifierFormatterBuilderRewriter {
	return &IdentifierFormatterBuilderRewriter{
		formatter:     formatter,
		formatterPass: NewIdentifierFormatterPass(formatter),
	}
}

func (rewriter *IdentifierFormatterBuilderRewriter) Rewrite(schemas Schemas, builders Builders) Builders {
	for i, builder := range builders {
		builders[i] = rewriter.rewriteBuilder(schemas, builder)
	}

	return builders
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteBuilder(schemas Schemas, builder Builder) Builder {
	builder.Package = rewriter.formatter.Package(builder.Package)
	builder.Name = rewriter.formatter.Object(builder.Name)
	builder.For = rewriter.formatterPass.processObject(schemas, builder.For)
	builder.Constructor = rewriter.rewriteConstructor(schemas, builder.Constructor)
	builder.Properties = tools.Map(builder.Properties, func(property StructField) StructField {
		structField := rewriter.formatterPass.processStructField(schemas, property)
		structField.Name = rewriter.formatter.ObjectPrivateField(structField.Name)
		return structField
	})
	builder.Options = tools.Map(builder.Options, func(option Option) Option {
		return rewriter.rewriteOption(schemas, option)
	})

	return builder
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteConstructor(schemas Schemas, constructor Constructor) Constructor {
	constructor.Args = tools.Map(constructor.Args, func(argument Argument) Argument {
		return rewriter.rewriteArgument(schemas, argument)
	})
	constructor.Assignments = tools.Map(constructor.Assignments, func(assignment Assignment) Assignment {
		return rewriter.rewriteAssignment(schemas, assignment)
	})

	return constructor
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteArgument(schemas Schemas, argument Argument) Argument {
	argument.Name = rewriter.formatter.Variable(argument.Name)
	argument.Type = rewriter.formatterPass.processType(schemas, argument.Type)

	return argument
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteAssignment(schemas Schemas, assignment Assignment) Assignment {
	assignment.Path = rewriter.rewritePath(schemas, assignment.Path)
	assignment.Value = rewriter.rewriteAssignmentValue(schemas, assignment.Value)
	assignment.Constraints = tools.Map(assignment.Constraints, func(constraint AssignmentConstraint) AssignmentConstraint {
		constraint.Argument = rewriter.rewriteArgument(schemas, constraint.Argument)
		return constraint
	})

	return assignment
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteAssignmentValue(schemas Schemas, value AssignmentValue) AssignmentValue {
	if value.Argument != nil {
		arg := rewriter.rewriteArgument(schemas, *value.Argument)
		value.Argument = &arg
	} else if value.Envelope != nil {
		value.Envelope.Type = rewriter.formatterPass.processType(schemas, value.Envelope.Type)
		value.Envelope.Values = tools.Map(value.Envelope.Values, func(envelopeField EnvelopeFieldValue) EnvelopeFieldValue {
			return rewriter.rewriteEnvelopeFieldValue(schemas, envelopeField)
		})
	}

	return value
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteEnvelopeFieldValue(schemas Schemas, fieldValue EnvelopeFieldValue) EnvelopeFieldValue {
	fieldValue.Path = rewriter.rewritePath(schemas, fieldValue.Path)
	fieldValue.Value = rewriter.rewriteAssignmentValue(schemas, fieldValue.Value)
	return fieldValue
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewritePath(schemas Schemas, path Path) Path {
	return tools.Map(path, func(item PathItem) PathItem {
		item.Identifier = rewriter.formatter.ObjectPublicField(item.Identifier)
		item.Type = rewriter.formatterPass.processType(schemas, item.Type)

		if item.TypeHint != nil {
			typeHint := *item.TypeHint
			typeHint = rewriter.formatterPass.processType(schemas, typeHint)
			item.TypeHint = &typeHint
		}

		return item
	})
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteOption(schemas Schemas, option Option) Option {
	option.Name = rewriter.formatter.Option(option.Name)
	option.Args = tools.Map(option.Args, func(arg Argument) Argument {
		return rewriter.rewriteArgument(schemas, arg)
	})
	option.Assignments = tools.Map(option.Assignments, func(assignment Assignment) Assignment {
		return rewriter.rewriteAssignment(schemas, assignment)
	})

	return option
}

// endregion
