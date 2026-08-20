use crate::cog;
use crate::types::sandbox;

#[derive(Debug, Clone)]
pub struct DashboardBuilder {
    internal: sandbox::Dashboard,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl DashboardBuilder {
    pub fn new() -> Self {
        Self {
            internal: sandbox::Dashboard::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for DashboardBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl DashboardBuilder {
    pub fn with_variable(mut self, name: String, value: String) -> Self {
        self.internal
            .variables
            .push(sandbox::Variable { name, value });

        self
    }
}

impl cog::Builder<sandbox::Dashboard> for DashboardBuilder {
    fn build(&self) -> Result<sandbox::Dashboard, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
