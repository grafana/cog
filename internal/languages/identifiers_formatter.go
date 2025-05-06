package languages

import (
	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/ast/compiler"
)

type IdentifiersConfig struct {
	PackageNameFunc func(string) string
	ObjectNameFunc  func(string) string
	FieldNameFunc   func(string) string
	BuilderNameFunc func(string) string
	OptionNameFunc  func(string) string
	ArgNameFunc     func(string) string
	AssignmentFunc  func(string) string
}

func FormatIdentifiers(language Language, context Context) (Context, error) {
	var err error
	identifiersConfig := IdentifiersConfig{}
	if identifierProvider, ok := language.(IdentifiersFormatter); ok {
		identifiersConfig = identifierProvider.Identifiers()
	}

	schemaVisitor := compiler.Visitor{
		OnSchema: func(visitor *compiler.Visitor, schema *ast.Schema) (*ast.Schema, error) {
			if identifiersConfig.PackageNameFunc != nil {
				schema.Package = identifiersConfig.PackageNameFunc(schema.Package)
			}
			return schema, nil
		},
		OnObject: func(visitor *compiler.Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
			if identifiersConfig.ObjectNameFunc != nil {
				object.Name = identifiersConfig.ObjectNameFunc(object.Name)
			}
			return object, nil
		},
		OnStructField: func(visitor *compiler.Visitor, schema *ast.Schema, field ast.StructField) (ast.StructField, error) {
			if identifiersConfig.FieldNameFunc != nil {
				field.Name = identifiersConfig.FieldNameFunc(field.Name)
			}
			return field, nil
		},
		OnRef: func(visitor *compiler.Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
			ref := def.AsRef()
			if identifiersConfig.PackageNameFunc != nil {
				ref.ReferredPkg = identifiersConfig.PackageNameFunc(ref.ReferredPkg)
			}
			if identifiersConfig.ObjectNameFunc != nil {
				ref.ReferredType = identifiersConfig.ObjectNameFunc(ref.ReferredType)
			}
			return ref.AsType(), nil
		},
		OnConstantRef: func(visitor *compiler.Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
			ref := def.AsConstantRef()
			if identifiersConfig.PackageNameFunc != nil {
				ref.ReferredPkg = identifiersConfig.PackageNameFunc(ref.ReferredPkg)
			}
			if identifiersConfig.ObjectNameFunc != nil {
				ref.ReferredType = identifiersConfig.ObjectNameFunc(ref.ReferredType)
			}
			return ref.AsType(), nil
		},
	}

	context.Schemas, err = schemaVisitor.VisitSchemas(context.Schemas)
	if err != nil {
		return Context{}, err
	}

	builderVisitor := ast.BuilderVisitor{
		OnBuilder: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder) (ast.Builder, error) {
			if identifiersConfig.BuilderNameFunc != nil {
				builder.Name = identifiersConfig.BuilderNameFunc(builder.Name)
			}
			return visitor.TraverseBuilder(schemas, builder)
		},
		OnOption: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, option ast.Option) (ast.Option, error) {
			if identifiersConfig.OptionNameFunc != nil {
				option.Name = identifiersConfig.OptionNameFunc(builder.Name)
			}
			return visitor.TraverseOption(schemas, builder, option)
		},
		OnArgument: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, argument ast.Argument) (ast.Argument, error) {
			if identifiersConfig.ArgNameFunc != nil {
				argument.Name = identifiersConfig.ArgNameFunc(builder.Name)
			}
			return argument, nil
		},
		OnAssignment: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, assignment ast.Assignment) (ast.Assignment, error) {
			return assignment, nil
		},
	}

	context.Builders, err = builderVisitor.Visit(context.Schemas, context.Builders)
	if err != nil {
		return Context{}, err
	}

	return context, nil
}
