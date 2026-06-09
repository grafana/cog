use crate::cog;
use crate::types::known_any::Config;
use crate::types::known_any::SomeStruct;

#[derive(Debug, Clone)]
pub struct SomeStructBuilder {
    internal: SomeStruct,
    errors: Vec<cog::BuildError>,
}

impl SomeStructBuilder {
    pub fn new() -> Self {
        Self {
            internal: SomeStruct::default(),
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
    pub fn title(mut self, title: String) -> Self {
        let config = Config { title };
        self.internal.config =
            Some(serde_json::to_value(config).expect("Config should serialize to JSON"));

        self
    }
}

impl cog::Builder<SomeStruct> for SomeStructBuilder {
    fn build(&self) -> Result<SomeStruct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
