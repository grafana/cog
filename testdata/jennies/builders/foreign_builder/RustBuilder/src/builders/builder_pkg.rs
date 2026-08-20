use crate::cog;
use crate::types::some_pkg;

#[derive(Debug, Clone)]
pub struct SomeNiceBuilderBuilder {
    internal: some_pkg::SomeStruct,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl SomeNiceBuilderBuilder {
    pub fn new() -> Self {
        Self {
            internal: some_pkg::SomeStruct::default(),
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

impl cog::Builder<some_pkg::SomeStruct> for SomeNiceBuilderBuilder {
    fn build(&self) -> Result<some_pkg::SomeStruct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
