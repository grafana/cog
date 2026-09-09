use crate::cog;
use crate::types::sandbox;

#[derive(Debug, Clone)]
pub struct SomeStructBuilder {
    internal: sandbox::SomeStruct,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl SomeStructBuilder {
    pub fn new(title: String) -> Self {
        let mut builder = Self {
            internal: sandbox::SomeStruct::default(),
            errors: Vec::new(),
        };
        builder.internal.title = title;

        builder
    }
}

impl SomeStructBuilder {
    pub fn title(mut self, title: String) -> Self {
        self.internal.title = title;

        self
    }
}

impl cog::Builder<sandbox::SomeStruct> for SomeStructBuilder {
    fn build(&self) -> Result<sandbox::SomeStruct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
