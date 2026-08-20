use crate::cog;
use crate::types::nullable_map_assignment;
use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct SomeStructBuilder {
    internal: nullable_map_assignment::SomeStruct,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl SomeStructBuilder {
    pub fn new() -> Self {
        Self {
            internal: nullable_map_assignment::SomeStruct::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for SomeStructBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl SomeStructBuilder {
    pub fn config(mut self, config: HashMap<String, String>) -> Self {
        self.internal.config = config;

        self
    }
}

impl cog::Builder<nullable_map_assignment::SomeStruct> for SomeStructBuilder {
    fn build(&self) -> Result<nullable_map_assignment::SomeStruct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
