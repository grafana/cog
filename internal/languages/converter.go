package languages

import (
	"fmt"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type MappingGuard struct {
	Path ast.Path

	NotNull bool
	// Or
	Op    ast.Op
	Value any
}

func (guard MappingGuard) String() string {
	if guard.NotNull {
		return fmt.Sprintf("%s != nil", guard.Path)
	}

	return fmt.Sprintf("%s %s %v", guard.Path, guard.Op, guard.Value)
}

type ArgumentMapping struct {
	ValuePath ast.Path     // direct mapping between a JSON value and an argument
	Builder   *ast.Builder //  the argument is built with a builder
}

func (argMapping ArgumentMapping) ValueType() ast.Type {
	return argMapping.ValuePath.Last().Type
}

type OptionMapping struct {
	Option ast.Option // option in the builder

	Repeat bool

	Guards []MappingGuard
	Args   []ArgumentMapping
}

type Converter struct {
	Package string

	Object  *ast.Object
	Builder *ast.Builder

	ConstructorArgs []ArgumentMapping

	Mappings []OptionMapping
}

type ConverterGenerator struct {
}

func (generator *ConverterGenerator) FromBuilder(context Context, builder ast.Builder) Converter {
	return Converter{
		Package: builder.Package,

		Object:  &builder.For,
		Builder: &builder,

		ConstructorArgs: generator.constructorArgs(builder),

		Mappings: tools.Map(builder.Options, func(option ast.Option) OptionMapping {
			return generator.convertOption(context, option)
		}),
	}
}

func (generator *ConverterGenerator) constructorArgs(builder ast.Builder) []ArgumentMapping {
	return tools.Map(builder.Constructor.Assignments, func(assignment ast.Assignment) ArgumentMapping {
		// "constructor options" are expected to only have a single assignment
		return ArgumentMapping{
			ValuePath: assignment.Path,
		}
	})
}

func (generator *ConverterGenerator) convertOption(context Context, option ast.Option) OptionMapping {
	repeat := false

	assignments := tools.Filter(option.Assignments, func(assignment ast.Assignment) bool {
		return assignment.Value.Constant == nil
	})

	mapping := OptionMapping{
		Option: option,
		Args: tools.Map(assignments, func(assignment ast.Assignment) ArgumentMapping {
			repeat = repeat || assignment.Method == ast.AppendAssignment

			builder, found := context.ResolveAsBuilder(assignment.Path.Last().Type)
			if !found {
				return ArgumentMapping{
					ValuePath: assignment.Path,
				}
			}

			return ArgumentMapping{
				ValuePath: assignment.Path,
				Builder:   &builder,
			}
		}),
	}

	mapping.Repeat = repeat && len(mapping.Args) == 1

	// conditions safeguarding the conversion of the current option
	guards := orderedmap.New[string, MappingGuard]()

	// TODO: define guards other than "not null" checks? (0, "", ...)
	// TODO: assignment method (direct vs append)
	// TODO: builders + array of builders (and array of array of builders, ...)
	// TODO: disjunctions
	// TODO: envelopes?
	for _, assignment := range option.Assignments {
		nullPathChunksGuards := generator.pathNotNullGuards(assignment.Path)
		for _, guard := range nullPathChunksGuards {
			guards.Set(guard.String(), guard)
		}

		if assignment.Value.Constant != nil {
			guard := MappingGuard{
				Path:  assignment.Path,
				Op:    ast.EqualOp,
				Value: assignment.Value.Constant,
			}
			guards.Set(guard.String(), guard)
			continue
		}

		// For arrays: ensure they're not empty
		if assignment.Path.Last().Type.IsArray() {
			guard := MappingGuard{
				Path:  assignment.Path,
				Op:    ast.MinLengthOp,
				Value: 1,
			}
			guards.Set(guard.String(), guard)
		}

		// For strings: ensure they're not empty
		if assignment.Path.Last().Type.IsScalar() && assignment.Path.Last().Type.AsScalar().ScalarKind == ast.KindString {
			guard := MappingGuard{
				Path:  assignment.Path,
				Op:    ast.NotEqualOp,
				Value: "",
			}
			guards.Set(guard.String(), guard)
		}

		// TODO: Envelope assignment?
		if assignment.Value.Envelope != nil {
			continue
		}
	}

	mapping.Guards = make([]MappingGuard, 0, guards.Len())
	guards.Iterate(func(_ string, guard MappingGuard) {
		mapping.Guards = append(mapping.Guards, guard)
	})

	return mapping
}

func (generator *ConverterGenerator) pathNotNullGuards(path ast.Path) []MappingGuard {
	var guards []MappingGuard

	for i, chunk := range path {
		chunkType := chunk.Type
		if chunk.TypeHint != nil {
			chunkType = *chunk.TypeHint
		}

		// TODO: this is language-specific
		maybeNull := chunkType.Nullable || chunkType.IsAnyOf(ast.KindMap, ast.KindArray)
		if !maybeNull {
			continue
		}

		guards = append(guards, MappingGuard{
			Path:    path[:i+1],
			NotNull: true,
		})
	}

	return guards
}
