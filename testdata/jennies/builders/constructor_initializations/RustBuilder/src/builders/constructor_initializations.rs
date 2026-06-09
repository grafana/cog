use crate::cog;
use crate::types::constructor_initializations::CursorMode;
use crate::types::constructor_initializations::SomePanel;

#[derive(Debug, Clone)]
pub struct SomePanelBuilder {
    internal: SomePanel,
}

impl SomePanelBuilder {
    pub fn new() -> Self {
        let mut builder = Self {
            internal: SomePanel::default(),
        };
        builder.internal.r#type = "panel_type".to_string();
        builder.internal.cursor = CursorMode::Tooltip;

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

impl cog::Builder<SomePanel> for SomePanelBuilder {
    fn build(&self) -> Result<SomePanel, Vec<cog::BuildError>> {
        Ok(self.internal.clone())
    }
}
