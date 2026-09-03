use crate::cog;
use crate::types::constructor_initializations;

#[derive(Debug, Clone)]
pub struct SomePanelBuilder {
    internal: constructor_initializations::SomePanel,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl SomePanelBuilder {
    pub fn new() -> Self {
        let mut builder = Self {
            internal: constructor_initializations::SomePanel::default(),
            errors: Vec::new(),
        };
        builder.internal.r#type = "panel_type".to_string();
        builder.internal.cursor = constructor_initializations::CursorMode::Tooltip;

        builder
    }
}

impl Default for SomePanelBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl SomePanelBuilder {
    pub fn title(mut self, title: String) -> Self {
        self.internal.title = title;

        self
    }
}

impl cog::Builder<constructor_initializations::SomePanel> for SomePanelBuilder {
    fn build(&self) -> Result<constructor_initializations::SomePanel, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
