use crate::cog;
use crate::types::map_of_disjunctions;
use std::collections::HashMap;

#[derive(Debug, Clone)]
pub struct ElementBuilder {
    internal: map_of_disjunctions::Element,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl ElementBuilder {
    pub fn new() -> Self {
        Self {
            internal: map_of_disjunctions::Element::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for ElementBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl ElementBuilder {
    pub fn panel(mut self, panel: impl cog::Builder<map_of_disjunctions::Panel>) -> Self {
        let built0 = match panel.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.panel = Some(built0);

        self
    }
}

impl ElementBuilder {
    pub fn library_panel(
        mut self,
        library_panel: impl cog::Builder<map_of_disjunctions::LibraryPanel>,
    ) -> Self {
        let built0 = match library_panel.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.library_panel = Some(built0);

        self
    }
}

impl cog::Builder<map_of_disjunctions::Element> for ElementBuilder {
    fn build(&self) -> Result<map_of_disjunctions::Element, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}

#[derive(Debug, Clone)]
pub struct PanelBuilder {
    internal: map_of_disjunctions::Panel,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl PanelBuilder {
    pub fn new() -> Self {
        let mut builder = Self {
            internal: map_of_disjunctions::Panel::default(),
            errors: Vec::new(),
        };
        builder.internal.kind = "Panel".to_string();

        builder
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

impl cog::Builder<map_of_disjunctions::Panel> for PanelBuilder {
    fn build(&self) -> Result<map_of_disjunctions::Panel, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}

#[derive(Debug, Clone)]
pub struct LibraryPanelBuilder {
    internal: map_of_disjunctions::LibraryPanel,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl LibraryPanelBuilder {
    pub fn new() -> Self {
        let mut builder = Self {
            internal: map_of_disjunctions::LibraryPanel::default(),
            errors: Vec::new(),
        };
        builder.internal.kind = "Library".to_string();

        builder
    }
}

impl Default for LibraryPanelBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl LibraryPanelBuilder {
    pub fn text(mut self, text: String) -> Self {
        self.internal.text = text;

        self
    }
}

impl cog::Builder<map_of_disjunctions::LibraryPanel> for LibraryPanelBuilder {
    fn build(&self) -> Result<map_of_disjunctions::LibraryPanel, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}

#[derive(Debug, Clone)]
pub struct DashboardBuilder {
    internal: map_of_disjunctions::Dashboard,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl DashboardBuilder {
    pub fn new() -> Self {
        Self {
            internal: map_of_disjunctions::Dashboard::default(),
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
    pub fn panels(
        mut self,
        panels: HashMap<String, impl cog::Builder<map_of_disjunctions::Element>>,
    ) -> Self {
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

impl cog::Builder<map_of_disjunctions::Dashboard> for DashboardBuilder {
    fn build(&self) -> Result<map_of_disjunctions::Dashboard, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}

#[derive(Debug, Clone)]
pub struct PanelOrLibraryPanelBuilder {
    internal: map_of_disjunctions::PanelOrLibraryPanel,
    pub(crate) errors: Vec<cog::BuildError>,
}

impl PanelOrLibraryPanelBuilder {
    pub fn new() -> Self {
        Self {
            internal: map_of_disjunctions::PanelOrLibraryPanel::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for PanelOrLibraryPanelBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl PanelOrLibraryPanelBuilder {
    pub fn panel(mut self, panel: impl cog::Builder<map_of_disjunctions::Panel>) -> Self {
        let built0 = match panel.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.panel = Some(built0);

        self
    }
}

impl PanelOrLibraryPanelBuilder {
    pub fn library_panel(
        mut self,
        library_panel: impl cog::Builder<map_of_disjunctions::LibraryPanel>,
    ) -> Self {
        let built0 = match library_panel.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.library_panel = Some(built0);

        self
    }
}

impl cog::Builder<map_of_disjunctions::PanelOrLibraryPanel> for PanelOrLibraryPanelBuilder {
    fn build(&self) -> Result<map_of_disjunctions::PanelOrLibraryPanel, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
