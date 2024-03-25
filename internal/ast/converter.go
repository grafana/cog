package ast

import (
	"fmt"

	"github.com/grafana/cog/internal/orderedmap"
	"github.com/grafana/cog/internal/tools"
)

type MappingGuard struct {
	Path Path

	NotNull bool
	// Or
	Op    Op
	Value any
}

func (guard MappingGuard) String() string {
	if guard.NotNull {
		return fmt.Sprintf("%s != nil", guard.Path)
	}

	return fmt.Sprintf("%s %s %v", guard.Path, guard.Op, guard.Value)
}

type OptionMapping struct {
	Option Option // option in the builder
	Paths  []Path // paths assigned by the option
	Guards []MappingGuard
}

type Converter struct {
	Package string

	Object  *Object
	Builder *Builder

	ConstructorArgs []Path

	Mappings []OptionMapping
}

type ConverterGenerator struct {
}

func (generator *ConverterGenerator) FromBuilder(builder Builder) Converter {
	return Converter{
		Package: builder.Package,

		Object:  &builder.For,
		Builder: &builder,

		ConstructorArgs: generator.constructorArgs(builder),

		Mappings: tools.Map(builder.Options, generator.convertOption),
	}
}

func (generator *ConverterGenerator) constructorArgs(builder Builder) []Path {
	constructorOpts := tools.Filter(builder.Options, func(option Option) bool {
		return option.IsConstructorArg
	})

	return tools.Map(constructorOpts, func(option Option) Path {
		// "constructor options" are expected to only have a single assignment
		return option.Assignments[0].Path
	})
}

func (generator *ConverterGenerator) convertOption(option Option) OptionMapping {
	mapping := OptionMapping{
		Option: option,
		Paths: tools.Map(option.Assignments, func(assignment Assignment) Path {
			return assignment.Path
		}),
	}

	guards := orderedmap.New[string, MappingGuard]()

	// TODO: define guards other than "not null" checks (0, "", ...)
	// TODO: assignment method (direct vs append)
	for _, assignment := range option.Assignments {
		nullPathChunksGuards := generator.pathNotNullGuards(assignment.Path)
		for _, guard := range nullPathChunksGuards {
			guards.Set(guard.String(), guard)
		}

		if assignment.Value.Constant != nil {
			guard := MappingGuard{
				Path:  assignment.Path,
				Op:    EqualOp,
				Value: assignment.Value.Constant,
			}
			guards.Set(guard.String(), guard)
			continue
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

func (generator *ConverterGenerator) pathNotNullGuards(path Path) []MappingGuard {
	var guards []MappingGuard

	for i, chunk := range path {
		chunkType := chunk.Type
		if chunk.TypeHint != nil {
			chunkType = *chunk.TypeHint
		}

		// TODO: this is language-specific
		maybeNull := chunkType.Nullable || chunkType.IsAnyOf(KindMap, KindArray)
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
