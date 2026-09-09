use crate::cog;
use crate::types::with_dashes;

#[derive(Debug, Clone)]
pub struct SomeNiceBuilderBuilder {
    internal: with_dashes::SomeStruct,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl SomeNiceBuilderBuilder {
    pub fn new() -> Self {
        Self {
            internal: with_dashes::SomeStruct::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for SomeNiceBuilderBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl SomeNiceBuilderBuilder {
    pub fn title(mut self, title: String) -> Self {
        self.internal.title = title;

        self
    }
}

impl cog::Builder<with_dashes::SomeStruct> for SomeNiceBuilderBuilder {
    fn build(&self) -> Result<with_dashes::SomeStruct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
