use crate::cog;
use crate::types::known_any;
use crate::types::known_any::Config;

#[derive(Debug, Clone)]
pub struct SomeStructBuilder {
    internal: known_any::SomeStruct,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl SomeStructBuilder {
    pub fn new() -> Self {
        Self {
            internal: known_any::SomeStruct::default(),
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
        let mut custom = self
            .internal
            .config
            .clone()
            .and_then(|raw| serde_json::from_value::<known_any::Config>(raw).ok())
            .unwrap_or_default();
        custom.title = title;
        self.internal.config =
            Some(serde_json::to_value(custom).expect("known_any::Config should serialize to JSON"));

        self
    }
}

impl cog::Builder<known_any::SomeStruct> for SomeStructBuilder {
    fn build(&self) -> Result<known_any::SomeStruct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
