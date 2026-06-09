use crate::cog;
use crate::types::discriminator_without_option::AnEnum;
use crate::types::discriminator_without_option::NoShowFieldOption;
use crate::types::discriminator_without_option::ShowFieldOption;

#[derive(Debug, Clone)]
pub struct NoShowFieldOptionBuilder {
    internal: NoShowFieldOption,
    errors: Vec<cog::BuildError>,
}

impl NoShowFieldOptionBuilder {
    pub fn new() -> Self {
        Self {
            internal: NoShowFieldOption::default(),
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

impl cog::Builder<NoShowFieldOption> for NoShowFieldOptionBuilder {
    fn build(&self) -> Result<NoShowFieldOption, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}

#[derive(Debug, Clone)]
pub struct ShowFieldOptionBuilder {
    internal: ShowFieldOption,
    errors: Vec<cog::BuildError>,
}

impl ShowFieldOptionBuilder {
    pub fn new() -> Self {
        Self {
            internal: ShowFieldOption::default(),
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
    pub fn field(mut self, field: AnEnum) -> Self {
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

impl cog::Builder<ShowFieldOption> for ShowFieldOptionBuilder {
    fn build(&self) -> Result<ShowFieldOption, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
