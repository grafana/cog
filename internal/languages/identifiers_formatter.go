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
				object.SelfRef.ReferredPkg = identifiersConfig.PackageNameFunc(schema.Package)
			}
			if identifiersConfig.ObjectNameFunc != nil {
				object.Name = identifiersConfig.ObjectNameFunc(object.Name)
				object.SelfRef.ReferredType = object.Name
			}
			return visitor.TransverseObject(schema, object)
		},
		OnStructField: func(visitor *compiler.Visitor, schema *ast.Schema, field ast.StructField) (ast.StructField, error) {
			if identifiersConfig.FieldNameFunc != nil {
				if field.OriginalName == "" {
					field.OriginalName = field.Name
				}
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

	// Update schema package names directly after the visitor runs.
	// The schema visitor creates a DeepCopy of the schema internally,
	// so modifications to schema.Package inside OnObject don't propagate.
	if identifiersConfig.PackageNameFunc != nil {
		for _, schema := range context.Schemas {
			schema.Package = identifiersConfig.PackageNameFunc(schema.Package)
		}
	}

	builderVisitor := ast.BuilderVisitor{
		OnBuilder: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder) (ast.Builder, error) {
			if identifiersConfig.PackageNameFunc != nil {
				builder.Package = identifiersConfig.PackageNameFunc(builder.Package)
			}
			if identifiersConfig.BuilderNameFunc != nil {
				builder.Name = identifiersConfig.BuilderNameFunc(builder.Name)
			}
			return visitor.TraverseBuilder(schemas, builder)
		},
		OnFactory: func(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, factory ast.BuilderFactory) (ast.BuilderFactory, error) {
			if identifiersConfig.OptionNameFunc != nil {
				factory.Name = identifiersConfig.OptionNameFunc(factory.Name)
			}
			if identifiersConfig.PackageNameFunc != nil {
				factory.BuilderRef.ReferredPkg = identifiersConfig.PackageNameFunc(factory.BuilderRef.ReferredPkg)
			}
			if identifiersConfig.ObjectNameFunc != nil {
				factory.BuilderRef.ReferredType = identifiersConfig.ObjectNameFunc(factory.BuilderRef.ReferredType)
			}
			for i, arg := range factory.Args {
				if identifiersConfig.ArgNameFunc != nil {
					factory.Args[i].Name = identifiersConfig.ArgNameFunc(arg.Name)
				}
			}
			for i, call := range factory.OptionCalls {
				if identifiersConfig.OptionNameFunc != nil {
					factory.OptionCalls[i].Name = identifiersConfig.OptionNameFunc(call.Name)
				}
				for j, param := range call.Parameters {
					if param.Argument != nil && identifiersConfig.ArgNameFunc != nil {
						factory.OptionCalls[i].Parameters[j].Argument.Name = identifiersConfig.ArgNameFunc(param.Argument.Name)
					}
					if param.Factory != nil && identifiersConfig.OptionNameFunc != nil {
						factory.OptionCalls[i].Parameters[j].Factory.Ref.Factory = identifiersConfig.OptionNameFunc(param.Factory.Ref.Factory)
					}
				}
			}
			return factory, nil
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
			if identifiersConfig.ArgNameFunc != nil {
				if assignment.Value.Argument != nil {
					assignment.Value.Argument.Name = identifiersConfig.ArgNameFunc(assignment.Value.Argument.Name)
				}
			}

			updateEnvelope(assignment.Value.Envelope, identifiersConfig.ArgNameFunc, identifiersConfig.AssignmentFunc)

			for i, nilCheck := range assignment.NilChecks {
				if identifiersConfig.AssignmentFunc != nil {
					for j, p := range nilCheck.Path {
						assignment.NilChecks[i].Path[j].Identifier = identifiersConfig.AssignmentFunc(p.Identifier)
					}
				}
				assignment.NilChecks[i].EmptyValueType = formatTypeRefs(nilCheck.EmptyValueType, identifiersConfig.PackageNameFunc, identifiersConfig.ObjectNameFunc)
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

// formatTypeRefs applies package/object name formatters to all refs within a type, recursively.
// This is needed for types stored outside the schema visitor's reach (e.g., NilCheck.EmptyValueType).
func formatTypeRefs(typeDef ast.Type, pkgFn func(string) string, objFn func(string) string) ast.Type {
	switch typeDef.Kind {
	case ast.KindRef:
		ref := typeDef.AsRef().DeepCopy()
		if pkgFn != nil {
			ref.ReferredPkg = pkgFn(ref.ReferredPkg)
		}
		if objFn != nil {
			ref.ReferredType = objFn(ref.ReferredType)
		}
		typeDef = typeDef.DeepCopy()
		typeDef.Ref = &ref
	case ast.KindArray:
		typeDef = typeDef.DeepCopy()
		valueType := formatTypeRefs(typeDef.AsArray().ValueType, pkgFn, objFn)
		typeDef.Array.ValueType = valueType
	case ast.KindMap:
		typeDef = typeDef.DeepCopy()
		valueType := formatTypeRefs(typeDef.AsMap().ValueType, pkgFn, objFn)
		typeDef.Map.ValueType = valueType
	}
	return typeDef
}

func updateEnvelope(envelope *ast.AssignmentEnvelope, argFn func(string) string, assignFn func(string) string) {
	if envelope == nil {
		return
	}

	for i, env := range envelope.Values {
		if env.Value.Argument != nil && argFn != nil {
			envelope.Values[i].Value.Argument.Name = argFn(env.Value.Argument.Name)
		}

		for j, p := range env.Path {
			if assignFn != nil {
				envelope.Values[i].Path[j].Identifier = assignFn(p.Identifier)
			}
		}

		if env.Value.Envelope != nil {
			updateEnvelope(envelope.Values[i].Value.Envelope, argFn, assignFn)
		}
	}
}
