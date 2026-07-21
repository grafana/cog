use crate::cog;
use crate::types::panelbuilder;

#[derive(Debug, Clone)]
pub struct PanelBuilder {
    internal: panelbuilder::Panel,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl PanelBuilder {
    pub fn new() -> Self {
        Self {
            internal: panelbuilder::Panel::default(),
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
    pub fn only_from_this_dashboard(mut self, only_from_this_dashboard: bool) -> Self {
        self.internal.only_from_this_dashboard = only_from_this_dashboard;

        self
    }
}

impl PanelBuilder {
    pub fn only_in_time_range(mut self, only_in_time_range: bool) -> Self {
        self.internal.only_in_time_range = only_in_time_range;

        self
    }
}

impl PanelBuilder {
    pub fn tags(mut self, tags: Vec<String>) -> Self {
        self.internal.tags = tags;

        self
    }
}

impl PanelBuilder {
    pub fn limit(mut self, limit: u32) -> Self {
        self.internal.limit = limit;

        self
    }
}

impl PanelBuilder {
    pub fn show_user(mut self, show_user: bool) -> Self {
        self.internal.show_user = show_user;

        self
    }
}

impl PanelBuilder {
    pub fn show_time(mut self, show_time: bool) -> Self {
        self.internal.show_time = show_time;

        self
    }
}

impl PanelBuilder {
    pub fn show_tags(mut self, show_tags: bool) -> Self {
        self.internal.show_tags = show_tags;

        self
    }
}

impl PanelBuilder {
    pub fn navigate_to_panel(mut self, navigate_to_panel: bool) -> Self {
        self.internal.navigate_to_panel = navigate_to_panel;

        self
    }
}

impl PanelBuilder {
    pub fn navigate_before(mut self, navigate_before: String) -> Self {
        self.internal.navigate_before = navigate_before;

        self
    }
}

impl PanelBuilder {
    pub fn navigate_after(mut self, navigate_after: String) -> Self {
        self.internal.navigate_after = navigate_after;

        self
    }
}

impl cog::Builder<panelbuilder::Panel> for PanelBuilder {
    fn build(&self) -> Result<panelbuilder::Panel, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
