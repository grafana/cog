use crate::cog;
use crate::types::other_pkg::Name;
use crate::types::some_pkg::Person;

#[derive(Debug, Clone)]
pub struct PersonBuilder {
    internal: Person,
    errors: Vec<cog::BuildError>,
}

impl PersonBuilder {
    pub fn new() -> Self {
        Self {
            internal: Person::default(),
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
    pub fn name(mut self, name: Name) -> Self {
        self.internal.name = name;

        self
    }
}

impl cog::Builder<Person> for PersonBuilder {
    fn build(&self) -> Result<Person, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
