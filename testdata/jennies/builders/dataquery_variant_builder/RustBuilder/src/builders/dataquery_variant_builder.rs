use crate::cog;
use crate::cog::variants;
use crate::types::dataquery_variant_builder;

#[derive(Debug, Clone)]
pub struct LokiBuilderBuilder {
    internal: dataquery_variant_builder::Loki,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl LokiBuilderBuilder {
    pub fn new() -> Self {
        Self {
            internal: dataquery_variant_builder::Loki::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for LokiBuilderBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl LokiBuilderBuilder {
    pub fn expr(mut self, expr: String) -> Self {
        self.internal.expr = expr;

        self
    }
}

impl cog::Builder<Box<dyn variants::Dataquery>> for LokiBuilderBuilder {
    fn build(&self) -> Result<Box<dyn variants::Dataquery>, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(Box::new(self.internal.clone()))
    }
}
