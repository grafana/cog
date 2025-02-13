import * as cog from '../cog';
import * as promql from '../promql';

export class FuncCallExprBuilder implements cog.Builder<promql.FuncCallExpr> {
    protected readonly internal: promql.FuncCallExpr;

    constructor() {
        this.internal = promql.defaultFuncCallExpr();
        this.internal.type = "funcCallExpr";
    }

    /**
     * Builds the object.
     */
    build(): promql.FuncCallExpr {
        return this.internal;
    }

    functionVal(functionVal: string): this {
        if (!(functionVal.length >= 1)) {
            throw new Error("functionVal.length must be >= 1");
        }
        this.internal.function = functionVal;
        return this;
    }

    args(args: cog.Builder<promql.Expr>[]): this {
        const argsResources = args.map(builder1 => builder1.build());
        this.internal.args = argsResources;
        return this;
    }

    // Modified by veneer 'Duplicate[args]'
    // Modified by veneer 'ArrayToAppend'
    arg(arg: cog.Builder<promql.Expr>): this {
        if (!this.internal.args) {
            this.internal.args = [];
        }
        const argResource = arg.build();
        this.internal.args.push(argResource);
        return this;
    }
}

/**
 * Returns the input vector with all sample values converted to their absolute value.
 * See https://prometheus.io/docs/prometheus/latest/querying/functions/#abs
 */
export function abs(v: cog.Builder<promql.Expr>): FuncCallExprBuilder {
	const builder = new FuncCallExprBuilder();
	builder.functionVal("abs");
	builder.arg(v);

	return builder;
}

/**
 * Returns an empty vector if the vector passed to it has any elements (floats or native histograms) and a 1-element vector with the value 1 if the vector passed to it has no elements.
 * This is useful for alerting on when no time series exist for a given metric name and label combination.
 * See https://prometheus.io/docs/prometheus/latest/querying/functions/#absent
 */
export function absent(v: cog.Builder<promql.Expr>): FuncCallExprBuilder {
	const builder = new FuncCallExprBuilder();
	builder.functionVal("absent");
	builder.arg(v);

	return builder;
}

/**
 * Returns pi.
 * See https://prometheus.io/docs/prometheus/latest/querying/functions/#trigonometric-functions
 */
export function pi(): FuncCallExprBuilder {
	const builder = new FuncCallExprBuilder();
	builder.functionVal("pi");

	return builder;
}

/**
 * Calculates the φ-quantile (0 ≤ φ ≤ 1) of the values in the specified interval.
 * See https://prometheus.io/docs/prometheus/latest/querying/functions/#aggregation_over_time
 */
export function quantileOverTime(phi: number,v: cog.Builder<promql.Expr>): FuncCallExprBuilder {
	const builder = new FuncCallExprBuilder();
	builder.functionVal("quantile_over_time");
	builder.arg(promql.n(phi));
	builder.arg(v);

	return builder;
}

