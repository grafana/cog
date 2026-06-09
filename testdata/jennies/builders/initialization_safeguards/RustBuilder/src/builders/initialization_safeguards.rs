use crate::cog;
use crate::types::initialization_safeguards::Options;
use crate::types::initialization_safeguards::SomePanel;

#[derive(Debug, Clone)]
pub struct SomePanelBuilder {
    internal: SomePanel,
    errors: Vec<cog::BuildError>,
}

impl SomePanelBuilder {
    pub fn new() -> Self {
        Self {
            internal: SomePanel::default(),
            errors: Vec::new(),
        }
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

impl SomePanelBuilder {
    pub fn show_legend(mut self, show: bool) -> Self {
        self.internal
            .options
            .get_or_insert_with(Options::default)
            .legend
            .show = show;

        self
    }
}

impl cog::Builder<SomePanel> for SomePanelBuilder {
    fn build(&self) -> Result<SomePanel, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
