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
		OnObject: func(visitor *compiler.Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
			if identifiersConfig.PackageNameFunc != nil {
				schema.Package = identifiersConfig.PackageNameFunc(schema.Package)
				object.SelfRef.ReferredPkg = identifiersConfig.PackageNameFunc(schema.Package)
			}
			if identifiersConfig.ObjectNameFunc != nil {
				object.Name = identifiersConfig.ObjectNameFunc(object.Name)
				object.SelfRef.ReferredType = identifiersConfig.ObjectNameFunc(object.Name)
			}
			return visitor.TransverseObject(schema, object)
		},
		OnStructField: func(visitor *compiler.Visitor, schema *ast.Schema, field ast.StructField) (ast.StructField, error) {
			if identifiersConfig.FieldNameFunc != nil {
				field.Name = identifiersConfig.FieldNameFunc(field.Name)
			}
			return visitor.TransverseStructField(schema, field)
		},
		OnRef: func(visitor *compiler.Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
			ref := def.AsRef().DeepCopy()
			if identifiersConfig.PackageNameFunc != nil {
				ref.ReferredPkg = identifiersConfig.PackageNameFunc(ref.ReferredPkg)
			}
			if identifiersConfig.ObjectNameFunc != nil {
				ref.ReferredType = identifiersConfig.ObjectNameFunc(ref.ReferredType)
			}
			return ast.Type{
				Kind:     ast.KindRef,
				Nullable: def.Nullable,
				Default:  def.Default,
				Ref:      &ref,
			}, nil
		},
		OnConstantRef: func(visitor *compiler.Visitor, schema *ast.Schema, def ast.Type) (ast.Type, error) {
			ref := def.AsConstantRef().DeepCopy()
			if identifiersConfig.PackageNameFunc != nil {
				ref.ReferredPkg = identifiersConfig.PackageNameFunc(ref.ReferredPkg)
			}
			if identifiersConfig.ObjectNameFunc != nil {
				ref.ReferredType = identifiersConfig.ObjectNameFunc(ref.ReferredType)
			}
			return ast.Type{
				Kind:              ast.KindConstantRef,
				Nullable:          def.Nullable,
				Default:           def.Default,
				ConstantReference: &ref,
			}, nil
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
				option.Name = identifiersConfig.OptionNameFunc(option.Name)
			}
			return visitor.TraverseOption(schemas, builder, option)
		},
		OnArgument: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, argument ast.Argument) (ast.Argument, error) {
			if identifiersConfig.ArgNameFunc != nil {
				argument.Name = identifiersConfig.ArgNameFunc(argument.Name)
			}
			return argument, nil
		},
		OnAssignment: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, assignment ast.Assignment) (ast.Assignment, error) {
			if identifiersConfig.AssignmentFunc != nil {
				for i, p := range assignment.Path {
					assignment.Path[i].Identifier = identifiersConfig.AssignmentFunc(p.Identifier)
				}
			}
			return assignment, nil
		},
	}

	context.Builders, err = builderVisitor.Visit(context.Schemas, context.Builders)
	if err != nil {
		return Context{}, err
	}

	return context, nil
}
