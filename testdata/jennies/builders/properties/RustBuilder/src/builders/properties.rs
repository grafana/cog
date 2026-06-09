use crate::cog;
use crate::types::properties::SomeStruct;

#[derive(Debug, Clone)]
pub struct SomeStructBuilder {
    internal: SomeStruct,
}

impl SomeStructBuilder {
    pub fn new() -> Self {
        Self {
            internal: SomeStruct::default(),
        }
    }
}

impl Default for SomeStructBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl SomeStructBuilder {
    pub fn id(mut self, id: i64) -> Self {
        self.internal.id = id;

        self
    }
}

impl cog::Builder<SomeStruct> for SomeStructBuilder {
    fn build(&self) -> Result<SomeStruct, Vec<cog::BuildError>> {
        Ok(self.internal.clone())
    }
}
