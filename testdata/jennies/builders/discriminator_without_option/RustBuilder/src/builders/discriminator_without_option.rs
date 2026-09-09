use crate::cog;
use crate::types::discriminator_without_option;

#[derive(Debug, Clone)]
pub struct NoShowFieldOptionBuilder {
    internal: discriminator_without_option::NoShowFieldOption,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl NoShowFieldOptionBuilder {
    pub fn new() -> Self {
        Self {
            internal: discriminator_without_option::NoShowFieldOption::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for NoShowFieldOptionBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl NoShowFieldOptionBuilder {
    pub fn text(mut self, text: String) -> Self {
        self.internal.text = text;

        self
    }
}

impl cog::Builder<discriminator_without_option::NoShowFieldOption> for NoShowFieldOptionBuilder {
    fn build(
        &self,
    ) -> Result<discriminator_without_option::NoShowFieldOption, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}

#[derive(Debug, Clone)]
pub struct ShowFieldOptionBuilder {
    internal: discriminator_without_option::ShowFieldOption,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl ShowFieldOptionBuilder {
    pub fn new() -> Self {
        Self {
            internal: discriminator_without_option::ShowFieldOption::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for ShowFieldOptionBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl ShowFieldOptionBuilder {
    pub fn field(mut self, field: discriminator_without_option::AnEnum) -> Self {
        self.internal.field = field;

        self
    }
}

impl ShowFieldOptionBuilder {
    pub fn text(mut self, text: String) -> Self {
        self.internal.text = text;

        self
    }
}

impl cog::Builder<discriminator_without_option::ShowFieldOption> for ShowFieldOptionBuilder {
    fn build(&self) -> Result<discriminator_without_option::ShowFieldOption, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
