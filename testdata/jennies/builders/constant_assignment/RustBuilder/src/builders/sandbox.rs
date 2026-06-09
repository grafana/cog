use crate::cog;
use crate::types::sandbox::SomeStruct;

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
    pub fn editable(mut self) -> Self {
        self.internal.editable = true;

        self
    }
}

impl SomeStructBuilder {
    pub fn readonly(mut self) -> Self {
        self.internal.editable = false;

        self
    }
}

impl SomeStructBuilder {
    pub fn auto_refresh(mut self) -> Self {
        self.internal.auto_refresh = Some(true);

        self
    }
}

impl SomeStructBuilder {
    pub fn no_auto_refresh(mut self) -> Self {
        self.internal.auto_refresh = Some(false);

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
