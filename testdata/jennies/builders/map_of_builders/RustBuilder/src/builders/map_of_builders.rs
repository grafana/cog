use crate::cog;
use crate::types::map_of_builders::Dashboard;
use crate::types::map_of_builders::Panel;
use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct PanelBuilder {
    internal: Panel,
    errors: Vec<cog::BuildError>,
}

impl PanelBuilder {
    pub fn new() -> Self {
        Self {
            internal: Panel::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for PanelBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl PanelBuilder {
    pub fn title(mut self, title: String) -> Self {
        self.internal.title = title;

        self
    }
}

impl cog::Builder<Panel> for PanelBuilder {
    fn build(&self) -> Result<Panel, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}

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
    pub fn panels(mut self, panels: HashMap<String, impl cog::Builder<Panel>>) -> Self {
        let mut built0 = std::collections::HashMap::new();
        for (key0, item0) in panels {
            let built1 = match item0.build() {
                Ok(val) => val,
                Err(mut err) => {
                    self.errors.append(&mut err);
                    return self;
                }
            };
            built0.insert(key0, built1);
        }
        self.internal.panels = built0;

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
