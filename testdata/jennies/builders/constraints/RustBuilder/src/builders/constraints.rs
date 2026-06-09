use crate::cog;
use crate::types::constraints::SomeStruct;

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
    pub fn id(mut self, id: u64) -> Self {
        if id < 5 {
            self.errors
                .push(cog::BuildError::new("id", "id must be >= 5".to_string()));
        }
        if id >= 10 {
            self.errors
                .push(cog::BuildError::new("id", "id must be < 10".to_string()));
        }
        self.internal.id = id;

        self
    }
}

impl SomeStructBuilder {
    pub fn title(mut self, title: String) -> Self {
        if title.chars().count() < 1 {
            self.errors.push(cog::BuildError::new(
                "title",
                "title length must be >= 1".to_string(),
            ));
        }
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
