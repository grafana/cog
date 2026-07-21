use crate::cog;
use crate::types::initialization_safeguards;

#[derive(Debug, Clone)]
pub struct SomePanelBuilder {
    internal: initialization_safeguards::SomePanel,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl SomePanelBuilder {
    pub fn new() -> Self {
        Self {
            internal: initialization_safeguards::SomePanel::default(),
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
            .get_or_insert_with(initialization_safeguards::Options::default)
            .legend
            .show = show;

        self
    }
}

impl cog::Builder<initialization_safeguards::SomePanel> for SomePanelBuilder {
    fn build(&self) -> Result<initialization_safeguards::SomePanel, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
