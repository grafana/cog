package languages

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
	"github.com/grafana/cog/internal/jennies/common"
	"github.com/grafana/cog/internal/tools"
)

func FormatIdentifiers(formatterProvider IdentifiersFormatterProvider, context common.Context) (common.Context, error) {
	var err error

	formatter := formatterProvider.IdentifiersFormatter()

	formatterPass := NewIdentifierFormatterPass(context.Schemas, formatter)
	context.Schemas, err = formatterPass.Process(context.Schemas)
	if err != nil {
		return context, err
	}

	buildersRewriter := NewIdentifierFormatterBuilderRewriter(context.Schemas, formatter)
	context.Builders, err = buildersRewriter.Rewrite(context.Schemas, context.Builders)
	if err != nil {
		return context, err
	}

	return context, nil
}

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
	visitor         *compiler.Visitor
	originalSchemas ast.Schemas
	formatter       *IdentifierFormatter
}

func NewIdentifierFormatterPass(originalSchemas ast.Schemas, formatter *IdentifierFormatter) *IdentifierFormatterPass {
	pass := &IdentifierFormatterPass{
		originalSchemas: originalSchemas,
		formatter:       formatter,
	}

	pass.visitor = &compiler.Visitor{
		OnSchema:      pass.processSchema,
		OnObject:      pass.processObject,
		OnEnum:        pass.processEnum,
		OnRef:         pass.processRef,
		OnStructField: pass.processStructField,
	}

	return pass
}

func (pass *IdentifierFormatterPass) Process(schemas []*ast.Schema) ([]*ast.Schema, error) {
	return pass.visitor.VisitSchemas(schemas)
}

func (pass *IdentifierFormatterPass) ProcessObject(object ast.Object) (ast.Object, error) {
	return pass.visitor.VisitObject(nil, object)
}

func (pass *IdentifierFormatterPass) ProcessType(def ast.Type) (ast.Type, error) {
	return pass.visitor.VisitType(nil, def)
}

func (pass *IdentifierFormatterPass) ProcessStructField(field ast.StructField) (ast.StructField, error) {
	return pass.visitor.VisitStructField(nil, field)
}

func (pass *IdentifierFormatterPass) processSchema(visitor *compiler.Visitor, schema *ast.Schema) (*ast.Schema, error) {
	schema.Package = pass.formatter.Package(schema.Package)
	schema.EntryPoint = pass.formatter.Object(schema.EntryPoint)

	var err error
	var obj ast.Object

	schema.Objects = schema.Objects.Map(func(_ string, object ast.Object) ast.Object {
		if err != nil {
			return object
		}

		obj, err = visitor.VisitObject(schema, object)

		return obj
	})

	return schema, err
}

func (pass *IdentifierFormatterPass) processObject(visitor *compiler.Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
	switch {
	case object.Type.IsEnum():
		object.Name = pass.formatter.EnumMember(object.Name)
	case object.Type.IsConcreteScalar():
		object.Name = pass.formatter.Constant(object.Name)
	default:
		object.Name = pass.formatter.Object(object.Name)
	}

	object.SelfRef.ReferredPkg = pass.formatter.Package(object.SelfRef.ReferredPkg)
	object.SelfRef.ReferredType = object.Name

	typeDef, err := visitor.VisitType(schema, object.Type)
	if err != nil {
		return object, err
	}

	object.Type = typeDef

	return object, nil
}

func (pass *IdentifierFormatterPass) processEnum(_ *compiler.Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	def.Enum.Values = tools.Map(def.Enum.Values, func(member ast.EnumValue) ast.EnumValue {
		member.Name = pass.formatter.EnumMember(member.Name)
		return member
	})

	return def, nil
}

func (pass *IdentifierFormatterPass) processRef(_ *compiler.Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	referredObj, found := pass.originalSchemas.LocateObject(def.Ref.ReferredPkg, def.Ref.ReferredType)
	if !found {
		return def, nil
	}

	def.Ref.ReferredPkg = pass.formatter.Package(def.Ref.ReferredPkg)

	switch {
	case referredObj.Type.IsEnum():
		def.Ref.ReferredType = pass.formatter.Enum(referredObj.Name)
	case referredObj.Type.IsConcreteScalar():
		def.Ref.ReferredType = pass.formatter.Constant(referredObj.Name)
	default:
		def.Ref.ReferredType = pass.formatter.Object(referredObj.Name)
	}

	return def, nil
}

func (pass *IdentifierFormatterPass) processStructField(visitor *compiler.Visitor, schema *ast.Schema, field ast.StructField) (ast.StructField, error) {
	field.Name = pass.formatter.ObjectPublicField(field.Name)

	typeDef, err := visitor.VisitType(schema, field.Type)
	if err != nil {
		return field, err
	}
	field.Type = typeDef

	return field, nil
}

// endregion

// region Identifier formatter builder rewriter

type IdentifierFormatterBuilderRewriter struct {
	formatter     *IdentifierFormatter
	formatterPass *IdentifierFormatterPass
}

func NewIdentifierFormatterBuilderRewriter(originalSchemas ast.Schemas, formatter *IdentifierFormatter) *IdentifierFormatterBuilderRewriter {
	return &IdentifierFormatterBuilderRewriter{
		formatter:     formatter,
		formatterPass: NewIdentifierFormatterPass(originalSchemas, formatter),
	}
}

func (rewriter *IdentifierFormatterBuilderRewriter) Rewrite(schemas ast.Schemas, builders ast.Builders) (ast.Builders, error) {
	visitor := &ast.BuilderVisitor{
		OnBuilder:    rewriter.rewriteBuilder,
		OnOption:     rewriter.rewriteOption,
		OnProperty:   rewriter.rewriteProperty,
		OnArgument:   rewriter.rewriteArgument,
		OnAssignment: rewriter.rewriteAssignment,
	}

	return visitor.Visit(schemas, builders)
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteBuilder(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder) (ast.Builder, error) {
	var err error

	builder.Package = rewriter.formatter.Package(builder.Package)
	builder.Name = rewriter.formatter.Object(builder.Name)

	forObj, err := rewriter.formatterPass.ProcessObject(builder.For)
	if err != nil {
		return builder, err
	}
	builder.For = forObj

	return visitor.TraverseBuilder(schemas, builder)
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteProperty(_ *ast.BuilderVisitor, _ ast.Schemas, _ ast.Builder, property ast.StructField) (ast.StructField, error) {
	var err error

	property, err = rewriter.formatterPass.ProcessStructField(property)
	if err != nil {
		return property, err
	}

	property.Name = rewriter.formatter.ObjectPrivateField(property.Name)

	return property, nil
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteArgument(_ *ast.BuilderVisitor, _ ast.Schemas, _ ast.Builder, argument ast.Argument) (ast.Argument, error) {
	var err error

	argument.Name = rewriter.formatter.Variable(argument.Name)
	argument.Type, err = rewriter.formatterPass.ProcessType(argument.Type)

	return argument, err
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteAssignment(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, assignment ast.Assignment) (ast.Assignment, error) {
	var arg ast.Argument
	var path ast.Path
	var err error

	path, err = rewriter.rewritePath(assignment.Path)
	if err != nil {
		return assignment, err
	}
	assignment.Path = path

	value, err := rewriter.rewriteAssignmentValue(visitor, builder, schemas, assignment.Value)
	if err != nil {
		return assignment, err
	}
	assignment.Value = value

	assignment.Constraints = tools.Map(assignment.Constraints, func(constraint ast.AssignmentConstraint) ast.AssignmentConstraint {
		if err != nil {
			return constraint
		}

		arg, err = visitor.VisitArgument(schemas, builder, constraint.Argument)
		constraint.Argument = arg

		return constraint
	})

	return assignment, err
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteAssignmentValue(visitor *ast.BuilderVisitor, builder ast.Builder, schemas ast.Schemas, value ast.AssignmentValue) (ast.AssignmentValue, error) {
	if value.Argument != nil {
		arg, err := visitor.VisitArgument(schemas, builder, *value.Argument)
		if err != nil {
			return value, err
		}

		value.Argument = &arg
	} else if value.Envelope != nil {
		typeDef, err := rewriter.formatterPass.ProcessType(value.Envelope.Type)
		if err != nil {
			return value, err
		}

		value.Envelope.Type = typeDef
		var fieldVal ast.EnvelopeFieldValue
		value.Envelope.Values = tools.Map(value.Envelope.Values, func(envelopeField ast.EnvelopeFieldValue) ast.EnvelopeFieldValue {
			if err != nil {
				return envelopeField
			}

			fieldVal, err = rewriter.rewriteEnvelopeFieldValue(visitor, builder, schemas, envelopeField)
			return fieldVal
		})
		if err != nil {
			return value, err
		}
	}

	return value, nil
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteEnvelopeFieldValue(visitor *ast.BuilderVisitor, builder ast.Builder, schemas ast.Schemas, fieldValue ast.EnvelopeFieldValue) (ast.EnvelopeFieldValue, error) {
	path, err := rewriter.rewritePath(fieldValue.Path)
	if err != nil {
		return fieldValue, err
	}
	fieldValue.Path = path

	value, err := rewriter.rewriteAssignmentValue(visitor, builder, schemas, fieldValue.Value)
	fieldValue.Value = value

	return fieldValue, err
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewritePath(path ast.Path) (ast.Path, error) {
	var err error
	var typeDef ast.Type
	var typeHintDef ast.Type
	path = tools.Map(path, func(item ast.PathItem) ast.PathItem {
		if err != nil {
			return item
		}

		item.Identifier = rewriter.formatter.ObjectPublicField(item.Identifier)
		typeDef, err = rewriter.formatterPass.ProcessType(item.Type)
		if err != nil {
			return item
		}

		item.Type = typeDef

		if item.TypeHint != nil {
			typeHintDef, err = rewriter.formatterPass.ProcessType(*item.TypeHint)
			if err != nil {
				return item
			}

			item.TypeHint = &typeHintDef
		}

		return item
	})

	return path, err
}

func (rewriter *IdentifierFormatterBuilderRewriter) rewriteOption(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, option ast.Option) (ast.Option, error) {
	option.Name = rewriter.formatter.Option(option.Name)

	return visitor.TraverseOption(schemas, builder, option)
}

// endregion
