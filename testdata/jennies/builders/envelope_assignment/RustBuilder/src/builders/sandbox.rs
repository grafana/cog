use crate::cog;
use crate::types::sandbox::Dashboard;
use crate::types::sandbox::Variable;

#[derive(Debug, Clone)]
pub struct DashboardBuilder {
    internal: Dashboard,
    errors: Vec<cog::BuildError>,
}

impl DashboardBuilder {
    pub fn new() -> Self {
        Self {
            internal: Dashboard::default(),
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
        self.internal.variables.push(Variable { name, value });

        self
    }
}

impl cog::Builder<Dashboard> for DashboardBuilder {
    fn build(&self) -> Result<Dashboard, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
