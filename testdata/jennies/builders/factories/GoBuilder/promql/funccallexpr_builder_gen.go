package promql

import (
	cog "github.com/grafana/cog/generated/cog"
)

var _ cog.Builder[FuncCallExpr] = (*FuncCallExprBuilder)(nil)

type FuncCallExprBuilder struct {
    internal *FuncCallExpr
    errors cog.BuildErrors
}

func NewFuncCallExprBuilder() *FuncCallExprBuilder {
	resource := NewFuncCallExpr()
	builder := &FuncCallExprBuilder{
		internal: resource,
		errors: make(cog.BuildErrors, 0),
	}
    builder.internal.Type = "funcCallExpr"

	return builder
}


	
// Returns the input vector with all sample values converted to their absolute value.
// See https://prometheus.io/docs/prometheus/latest/querying/functions/#abs
func Abs(v cog.Builder[Expr]) *FuncCallExprBuilder {
	builder := NewFuncCallExprBuilder()
    builder.Function("abs")
    builder.Arg(v)

	return builder
}
	
// Returns an empty vector if the vector passed to it has any elements (floats or native histograms) and a 1-element vector with the value 1 if the vector passed to it has no elements.
// This is useful for alerting on when no time series exist for a given metric name and label combination.
// See https://prometheus.io/docs/prometheus/latest/querying/functions/#absent
func Absent(v cog.Builder[Expr]) *FuncCallExprBuilder {
	builder := NewFuncCallExprBuilder()
    builder.Function("absent")
    builder.Arg(v)

	return builder
}
	
// Returns pi.
// See https://prometheus.io/docs/prometheus/latest/querying/functions/#trigonometric-functions
func Pi() *FuncCallExprBuilder {
	builder := NewFuncCallExprBuilder()
    builder.Function("pi")

	return builder
}
	
// Calculates the φ-quantile (0 ≤ φ ≤ 1) of the values in the specified interval.
// See https://prometheus.io/docs/prometheus/latest/querying/functions/#aggregation_over_time
func QuantileOverTime(phi float64,v cog.Builder[Expr]) *FuncCallExprBuilder {
	builder := NewFuncCallExprBuilder()
    builder.Function("quantile_over_time")
    builder.Arg(N(phi))
    builder.Arg(v)

	return builder
}

func (builder *FuncCallExprBuilder) Build() (FuncCallExpr, error) {
	if err := builder.internal.Validate(); err != nil {
		return FuncCallExpr{}, err
	}
	
	if len(builder.errors) > 0 {
	    return FuncCallExpr{}, cog.MakeBuildErrors("promql.funcCallExpr", builder.errors)
	}

	return *builder.internal, nil
}

func (builder *FuncCallExprBuilder) Function(function string) *FuncCallExprBuilder {
    builder.internal.Function = function

    return builder
}

func (builder *FuncCallExprBuilder) Args(args []cog.Builder[Expr]) *FuncCallExprBuilder {
        argsResources := make([]Expr, 0, len(args))
        for _, r1 := range args {
                argsDepth1, err := r1.Build()
                if err != nil {
                    builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
                    return builder
                }
                argsResources = append(argsResources, argsDepth1)
        }
    builder.internal.Args = argsResources

    return builder
}

// Modified by veneer 'Duplicate[args]'
// Modified by veneer 'ArrayToAppend'
func (builder *FuncCallExprBuilder) Arg(arg cog.Builder[Expr]) *FuncCallExprBuilder {
    argResource, err := arg.Build()
    if err != nil {
        builder.errors = append(builder.errors, err.(cog.BuildErrors)...)
        return builder
    }
    builder.internal.Args = append(builder.internal.Args, argResource)

    return builder
}

