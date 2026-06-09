use crate::cog;
use crate::types::some_pkg::SomeStruct;

#[derive(Debug, Clone)]
pub struct SomeNiceBuilderBuilder {
    internal: SomeStruct,
    errors: Vec<cog::BuildError>,
}

impl SomeNiceBuilderBuilder {
    pub fn new() -> Self {
        Self {
            internal: SomeStruct::default(),
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

impl cog::Builder<SomeStruct> for SomeNiceBuilderBuilder {
    fn build(&self) -> Result<SomeStruct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
