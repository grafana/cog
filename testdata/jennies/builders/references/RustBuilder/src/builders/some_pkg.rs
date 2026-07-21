use crate::cog;
use crate::types::other_pkg;
use crate::types::some_pkg;

#[derive(Debug, Clone)]
pub struct PersonBuilder {
    internal: some_pkg::Person,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl PersonBuilder {
    pub fn new() -> Self {
        Self {
            internal: some_pkg::Person::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for PersonBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl PersonBuilder {
    pub fn name(mut self, name: other_pkg::Name) -> Self {
        self.internal.name = name;

        self
    }
}

impl cog::Builder<some_pkg::Person> for PersonBuilder {
    fn build(&self) -> Result<some_pkg::Person, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
