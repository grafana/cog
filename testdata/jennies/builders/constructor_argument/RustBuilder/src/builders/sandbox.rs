use crate::cog;
use crate::types::sandbox::SomeStruct;

#[derive(Debug, Clone)]
pub struct SomeStructBuilder {
    internal: SomeStruct,
    errors: Vec<cog::BuildError>,
}

impl SomeStructBuilder {
    pub fn new(title: String) -> Self {
        let mut builder = Self {
            internal: SomeStruct::default(),
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

impl cog::Builder<SomeStruct> for SomeStructBuilder {
    fn build(&self) -> Result<SomeStruct, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
