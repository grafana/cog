import typing
from ..cog import builder as cogbuilder
from ..models import promql


class FuncCallExpr(cogbuilder.Builder[promql.FuncCallExpr]):
    _internal: promql.FuncCallExpr

    def __init__(self):
        self._internal = promql.FuncCallExpr()        
        self._internal.type_val = "funcCallExpr"

    def build(self) -> promql.FuncCallExpr:
        """
        Builds the object.
        """
        return self._internal    
    
    def function(self, function: str) -> typing.Self:    
        if not len(function) >= 1:
            raise ValueError("len(function) must be >= 1")
        self._internal.function = function
    
        return self
    
    def args(self, args: list[cogbuilder.Builder[promql.Expr]]) -> typing.Self:    
        args_resources = [r1.build() for r1 in args]
        self._internal.args = args_resources
    
        return self
    
    def arg(self, arg: cogbuilder.Builder[promql.Expr]) -> typing.Self:    
        """
        Modified by veneer 'Duplicate[args]'
        Modified by veneer 'ArrayToAppend'
        """
            
        if self._internal.args is None:
            self._internal.args = []
        
        arg_resource = arg.build()
        self._internal.args.append(arg_resource)
    
        return self
    

"""
Returns the input vector with all sample values converted to their absolute value.
See https://prometheus.io/docs/prometheus/latest/querying/functions/#abs
"""
def abs(v: cogbuilder.Builder[promql.Expr]):
    builder = FuncCallExpr()
    builder.function("abs")
    builder.arg(v)

    return builder

"""
Returns an empty vector if the vector passed to it has any elements (floats or native histograms) and a 1-element vector with the value 1 if the vector passed to it has no elements.
This is useful for alerting on when no time series exist for a given metric name and label combination.
See https://prometheus.io/docs/prometheus/latest/querying/functions/#absent
"""
def absent(v: cogbuilder.Builder[promql.Expr]):
    builder = FuncCallExpr()
    builder.function("absent")
    builder.arg(v)

    return builder

"""
Returns pi.
See https://prometheus.io/docs/prometheus/latest/querying/functions/#trigonometric-functions
"""
def pi():
    builder = FuncCallExpr()
    builder.function("pi")

    return builder

"""
Calculates the φ-quantile (0 ≤ φ ≤ 1) of the values in the specified interval.
See https://prometheus.io/docs/prometheus/latest/querying/functions/#aggregation_over_time
"""
def quantile_over_time(phi: float, v: cogbuilder.Builder[promql.Expr]):
    builder = FuncCallExpr()
    builder.function("quantile_over_time")
    builder.arg(n(phi))
    builder.arg(v)

    return builder
