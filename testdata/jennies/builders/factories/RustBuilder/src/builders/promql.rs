use crate::cog;
use crate::types::promql::Expr;
use crate::types::promql::FuncCallExpr;

#[derive(Debug, Clone)]
pub struct FuncCallExprBuilder {
    internal: FuncCallExpr,
    errors: Vec<cog::BuildError>,
}

impl FuncCallExprBuilder {
    pub fn new() -> Self {
        let mut builder = Self {
            internal: FuncCallExpr::default(),
            errors: Vec::new(),
        };
        builder.internal.r#type = "funcCallExpr".to_string();

        builder
    }
}

impl Default for FuncCallExprBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl FuncCallExprBuilder {
    pub fn function(mut self, function: String) -> Self {
        if function.chars().count() < 1 {
            self.errors.push(cog::BuildError::new(
                "function",
                "function length must be >= 1".to_string(),
            ));
        }
        self.internal.function = function;

        self
    }
}

impl FuncCallExprBuilder {
    pub fn args(mut self, args: Vec<impl cog::Builder<Expr>>) -> Self {
        let mut built0 = Vec::new();
        for item0 in args {
            let built1 = match item0.build() {
                Ok(val) => val,
                Err(mut err) => {
                    self.errors.append(&mut err);
                    return self;
                }
            };
            built0.push(built1);
        }
        self.internal.args = built0;

        self
    }
}

/// Modified by veneer 'Duplicate[args]'
/// Modified by veneer 'ArrayToAppend'
impl FuncCallExprBuilder {
    pub fn arg(mut self, arg: impl cog::Builder<Expr>) -> Self {
        let built0 = match arg.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.args.push(built0);

        self
    }
}

/// Returns the input vector with all sample values converted to their absolute value.
/// See https://prometheus.io/docs/prometheus/latest/querying/functions/#abs
impl FuncCallExprBuilder {
    pub fn abs(v: impl cog::Builder<Expr>) -> Self {
        Self::new().function("abs".to_string()).arg(v)
    }
}

/// Returns an empty vector if the vector passed to it has any elements (floats or native histograms) and a 1-element vector with the value 1 if the vector passed to it has no elements.
/// This is useful for alerting on when no time series exist for a given metric name and label combination.
/// See https://prometheus.io/docs/prometheus/latest/querying/functions/#absent
impl FuncCallExprBuilder {
    pub fn absent(v: impl cog::Builder<Expr>) -> Self {
        Self::new().function("absent".to_string()).arg(v)
    }
}

/// Returns pi.
/// See https://prometheus.io/docs/prometheus/latest/querying/functions/#trigonometric-functions
impl FuncCallExprBuilder {
    pub fn pi() -> Self {
        Self::new().function("pi".to_string())
    }
}

/// Calculates the φ-quantile (0 ≤ φ ≤ 1) of the values in the specified interval.
/// See https://prometheus.io/docs/prometheus/latest/querying/functions/#aggregation_over_time
impl FuncCallExprBuilder {
    pub fn quantile_over_time(phi: f64, v: impl cog::Builder<Expr>) -> Self {
        Self::new()
            .function("quantile_over_time".to_string())
            .arg(NumberLiteralExprBuilder::n(phi))
            .arg(v)
    }
}

impl cog::Builder<FuncCallExpr> for FuncCallExprBuilder {
    fn build(&self) -> Result<FuncCallExpr, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
