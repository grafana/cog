package ast

import (
	"github.com/grafana/cog/internal/tools"
)

type VisitBuilderFunc func(visitor *BuilderVisitor, schemas Schemas, builder Builder) Builder
type VisitPropertyFunc func(visitor *BuilderVisitor, schemas Schemas, builder Builder, property StructField) StructField
type VisitConstructorFunc func(visitor *BuilderVisitor, schemas Schemas, builder Builder, constructor Constructor) Constructor
type VisitOptionFunc func(visitor *BuilderVisitor, schemas Schemas, builder Builder, option Option) Option
type VisitArgumentFunc func(visitor *BuilderVisitor, schemas Schemas, builder Builder, argument Argument) Argument
type VisitAssignmentFunc func(visitor *BuilderVisitor, schemas Schemas, builder Builder, assignment Assignment) Assignment

type BuilderVisitor struct {
	OnBuilder     VisitBuilderFunc
	OnProperty    VisitPropertyFunc
	OnConstructor VisitConstructorFunc
	OnOption      VisitOptionFunc
	OnArgument    VisitArgumentFunc
	OnAssignment  VisitAssignmentFunc
}

func (visitor *BuilderVisitor) Visit(schemas Schemas, builders Builders) Builders {
	for i, builder := range builders {
		builders[i] = visitor.visitBuilder(schemas, builder)
	}

	return builders
}

func (visitor *BuilderVisitor) visitBuilder(schemas Schemas, builder Builder) Builder {
	if visitor.OnBuilder != nil {
		return visitor.OnBuilder(visitor, schemas, builder)
	}

	builder.Constructor = visitor.visitConstructor(schemas, builder, builder.Constructor)
	builder.Properties = tools.Map(builder.Properties, func(property StructField) StructField {
		return visitor.visitProperty(schemas, builder, property)
	})
	builder.Options = tools.Map(builder.Options, func(option Option) Option {
		return visitor.visitOption(schemas, builder, option)
	})

	return builder
}

func (visitor *BuilderVisitor) visitProperty(schemas Schemas, builder Builder, property StructField) StructField {
	if visitor.OnProperty != nil {
		return visitor.OnProperty(visitor, schemas, builder, property)
	}

	return property
}

func (visitor *BuilderVisitor) visitConstructor(schemas Schemas, builder Builder, constructor Constructor) Constructor {
	if visitor.OnConstructor != nil {
		return visitor.OnConstructor(visitor, schemas, builder, constructor)
	}
	constructor.Args = tools.Map(constructor.Args, func(argument Argument) Argument {
		return visitor.visitArgument(schemas, builder, argument)
	})
	constructor.Assignments = tools.Map(constructor.Assignments, func(assignment Assignment) Assignment {
		return visitor.visitAssignment(schemas, builder, assignment)
	})

	return constructor
}

func (visitor *BuilderVisitor) visitArgument(schemas Schemas, builder Builder, argument Argument) Argument {
	if visitor.OnArgument != nil {
		return visitor.OnArgument(visitor, schemas, builder, argument)
	}

	return argument
}

func (visitor *BuilderVisitor) visitAssignment(schemas Schemas, builder Builder, assignment Assignment) Assignment {
	if visitor.OnAssignment != nil {
		return visitor.OnAssignment(visitor, schemas, builder, assignment)
	}

	assignment.Constraints = tools.Map(assignment.Constraints, func(constraint AssignmentConstraint) AssignmentConstraint {
		constraint.Argument = visitor.visitArgument(schemas, builder, constraint.Argument)
		return constraint
	})

	return assignment
}

func (visitor *BuilderVisitor) visitOption(schemas Schemas, builder Builder, option Option) Option {
	if visitor.OnOption != nil {
		return visitor.OnOption(visitor, schemas, builder, option)
	}

	option.Args = tools.Map(option.Args, func(arg Argument) Argument {
		return visitor.visitArgument(schemas, builder, arg)
	})
	option.Assignments = tools.Map(option.Assignments, func(assignment Assignment) Assignment {
		return visitor.visitAssignment(schemas, builder, assignment)
	})

	return option
}
