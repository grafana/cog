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

	schemaVisitor := identifiersConfig.schemaVisitor()
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
	if identifiersConfig.ObjectNameFunc != nil {
		for _, schema := range context.Schemas {
			if schema.EntryPoint != "" {
				schema.EntryPoint = identifiersConfig.ObjectNameFunc(schema.EntryPoint)
			}
		}
	}

	builderVisitor := identifiersConfig.builderVisitor()
	context.Builders, err = builderVisitor.Visit(context.Schemas, context.Builders)
	if err != nil {
		return Context{}, err
	}

	return context, nil
}

func (cfg IdentifiersConfig) schemaVisitor() compiler.Visitor {
	return compiler.Visitor{
		OnObject:      cfg.onObject,
		OnStructField: cfg.onStructField,
		OnRef:         cfg.onRef,
		OnEnum:        cfg.onEnum,
		OnConstantRef: cfg.onConstantRef,
	}
}

func (cfg IdentifiersConfig) builderVisitor() ast.BuilderVisitor {
	return ast.BuilderVisitor{
		OnBuilder:    cfg.onBuilder,
		OnFactory:    cfg.onFactory,
		OnOption:     cfg.onOption,
		OnArgument:   cfg.onArgument,
		OnAssignment: cfg.onAssignment,
	}
}

func (cfg IdentifiersConfig) onObject(visitor *compiler.Visitor, schema *ast.Schema, object ast.Object) (ast.Object, error) {
	if cfg.PackageNameFunc != nil {
		object.SelfRef.ReferredPkg = cfg.PackageNameFunc(schema.Package)
	}
	if cfg.ObjectNameFunc != nil {
		object.Name = cfg.ObjectNameFunc(object.Name)
		object.SelfRef.ReferredType = cfg.ObjectNameFunc(object.SelfRef.ReferredType)
	}
	return visitor.TransverseObject(schema, object)
}

func (cfg IdentifiersConfig) onStructField(visitor *compiler.Visitor, schema *ast.Schema, field ast.StructField) (ast.StructField, error) {
	if cfg.FieldNameFunc != nil {
		if field.OriginalName == "" {
			field.OriginalName = field.Name
		}
		field.Name = cfg.FieldNameFunc(field.Name)
	}
	return visitor.TransverseStructField(schema, field)
}

func (cfg IdentifiersConfig) onRef(_ *compiler.Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	ref := def.AsRef().DeepCopy()
	if cfg.PackageNameFunc != nil {
		ref.ReferredPkg = cfg.PackageNameFunc(ref.ReferredPkg)
	}
	if cfg.ObjectNameFunc != nil {
		ref.ReferredType = cfg.ObjectNameFunc(ref.ReferredType)
	}
	return ast.Type{
		Kind:     ast.KindRef,
		Nullable: def.Nullable,
		Default:  def.Default,
		Ref:      &ref,
	}, nil
}

func (cfg IdentifiersConfig) onEnum(_ *compiler.Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	if cfg.ObjectNameFunc == nil {
		return def, nil
	}
	enumType := def.AsEnum().DeepCopy()
	for i, val := range enumType.Values {
		enumType.Values[i].Name = cfg.ObjectNameFunc(val.Name)
	}
	def.Enum = &enumType
	return def, nil
}

func (cfg IdentifiersConfig) onConstantRef(_ *compiler.Visitor, _ *ast.Schema, def ast.Type) (ast.Type, error) {
	ref := def.AsConstantRef().DeepCopy()
	if cfg.PackageNameFunc != nil {
		ref.ReferredPkg = cfg.PackageNameFunc(ref.ReferredPkg)
	}
	if cfg.ObjectNameFunc != nil {
		ref.ReferredType = cfg.ObjectNameFunc(ref.ReferredType)
	}
	return ast.Type{
		Kind:              ast.KindConstantRef,
		Nullable:          def.Nullable,
		Default:           def.Default,
		ConstantReference: &ref,
	}, nil
}

func (cfg IdentifiersConfig) onBuilder(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder) (ast.Builder, error) {
	if cfg.PackageNameFunc != nil {
		if builder.OriginalPackage == "" {
			builder.OriginalPackage = builder.Package
		}
		builder.Package = cfg.PackageNameFunc(builder.Package)
		builder.For.SelfRef.ReferredPkg = cfg.PackageNameFunc(builder.For.SelfRef.ReferredPkg)
	}
	if cfg.ObjectNameFunc != nil {
		builder.For.Name = cfg.ObjectNameFunc(builder.For.Name)
		builder.For.SelfRef.ReferredType = builder.For.Name
	}
	if cfg.BuilderNameFunc != nil {
		builder.Name = cfg.BuilderNameFunc(builder.Name)
	}
	return visitor.TraverseBuilder(schemas, builder)
}

func (cfg IdentifiersConfig) onFactory(_ *ast.BuilderVisitor, _ ast.Schemas, _ ast.Builder, factory ast.BuilderFactory) (ast.BuilderFactory, error) {
	if cfg.OptionNameFunc != nil {
		factory.Name = cfg.OptionNameFunc(factory.Name)
	}
	if cfg.PackageNameFunc != nil {
		factory.BuilderRef.ReferredPkg = cfg.PackageNameFunc(factory.BuilderRef.ReferredPkg)
	}
	if cfg.ObjectNameFunc != nil {
		factory.BuilderRef.ReferredType = cfg.ObjectNameFunc(factory.BuilderRef.ReferredType)
	}
	for i, arg := range factory.Args {
		if cfg.ArgNameFunc != nil {
			factory.Args[i].Name = cfg.ArgNameFunc(arg.Name)
		}
	}
	for i, call := range factory.OptionCalls {
		if cfg.OptionNameFunc != nil {
			factory.OptionCalls[i].Name = cfg.OptionNameFunc(call.Name)
		}
		for j, param := range call.Parameters {
			if param.Argument != nil && cfg.ArgNameFunc != nil {
				factory.OptionCalls[i].Parameters[j].Argument.Name = cfg.ArgNameFunc(param.Argument.Name)
			}
			if param.Factory != nil && cfg.OptionNameFunc != nil {
				factory.OptionCalls[i].Parameters[j].Factory.Ref.Factory = cfg.OptionNameFunc(param.Factory.Ref.Factory)
			}
		}
	}
	return factory, nil
}

func (cfg IdentifiersConfig) onOption(visitor *ast.BuilderVisitor, schemas ast.Schemas, builder ast.Builder, option ast.Option) (ast.Option, error) {
	if cfg.OptionNameFunc != nil {
		option.Name = cfg.OptionNameFunc(option.Name)
	}
	return visitor.TraverseOption(schemas, builder, option)
}

func (cfg IdentifiersConfig) onArgument(_ *ast.BuilderVisitor, _ ast.Schemas, _ ast.Builder, argument ast.Argument) (ast.Argument, error) {
	if cfg.ArgNameFunc != nil {
		argument.Name = cfg.ArgNameFunc(argument.Name)
	}
	argument.Type = formatTypeRefs(argument.Type, cfg.PackageNameFunc, cfg.ObjectNameFunc)
	return argument, nil
}

func (cfg IdentifiersConfig) onAssignment(_ *ast.BuilderVisitor, _ ast.Schemas, _ ast.Builder, assignment ast.Assignment) (ast.Assignment, error) {
	for i, p := range assignment.Path {
		if cfg.AssignmentFunc != nil {
			assignment.Path[i].Identifier = cfg.AssignmentFunc(p.Identifier)
		}
		assignment.Path[i].Type = formatTypeRefs(p.Type, cfg.PackageNameFunc, cfg.ObjectNameFunc)
		if p.TypeHint != nil {
			hint := formatTypeRefs(*p.TypeHint, cfg.PackageNameFunc, cfg.ObjectNameFunc)
			assignment.Path[i].TypeHint = &hint
		}
	}
	if assignment.Value.Argument != nil {
		if cfg.ArgNameFunc != nil {
			assignment.Value.Argument.Name = cfg.ArgNameFunc(assignment.Value.Argument.Name)
		}
		assignment.Value.Argument.Type = formatTypeRefs(assignment.Value.Argument.Type, cfg.PackageNameFunc, cfg.ObjectNameFunc)
	}

	updateEnvelope(assignment.Value.Envelope, cfg.ArgNameFunc, cfg.AssignmentFunc, cfg.PackageNameFunc, cfg.ObjectNameFunc)

	for i, nilCheck := range assignment.NilChecks {
		for j, p := range nilCheck.Path {
			if cfg.AssignmentFunc != nil {
				assignment.NilChecks[i].Path[j].Identifier = cfg.AssignmentFunc(p.Identifier)
			}
			assignment.NilChecks[i].Path[j].Type = formatTypeRefs(p.Type, cfg.PackageNameFunc, cfg.ObjectNameFunc)
			if p.TypeHint != nil {
				hint := formatTypeRefs(*p.TypeHint, cfg.PackageNameFunc, cfg.ObjectNameFunc)
				assignment.NilChecks[i].Path[j].TypeHint = &hint
			}
		}
		assignment.NilChecks[i].EmptyValueType = formatTypeRefs(nilCheck.EmptyValueType, cfg.PackageNameFunc, cfg.ObjectNameFunc)
	}

	return assignment, nil
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
	case ast.KindDisjunction:
		typeDef = typeDef.DeepCopy()
		for i, branch := range typeDef.Disjunction.Branches {
			typeDef.Disjunction.Branches[i] = formatTypeRefs(branch, pkgFn, objFn)
		}
	}
	return typeDef
}

func updateEnvelope(envelope *ast.AssignmentEnvelope, argFn func(string) string, assignFn func(string) string, pkgFn func(string) string, objFn func(string) string) {
	if envelope == nil {
		return
	}

	for i, env := range envelope.Values {
		if env.Value.Argument != nil {
			if argFn != nil {
				envelope.Values[i].Value.Argument.Name = argFn(env.Value.Argument.Name)
			}
			envelope.Values[i].Value.Argument.Type = formatTypeRefs(env.Value.Argument.Type, pkgFn, objFn)
		}

		for j, p := range env.Path {
			if assignFn != nil {
				envelope.Values[i].Path[j].Identifier = assignFn(p.Identifier)
			}
			envelope.Values[i].Path[j].Type = formatTypeRefs(p.Type, pkgFn, objFn)
		}

		if env.Value.Envelope != nil {
			updateEnvelope(envelope.Values[i].Value.Envelope, argFn, assignFn, pkgFn, objFn)
		}
	}
}
