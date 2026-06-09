use crate::cog;
use crate::types::builder_delegation::Dashboard;
use crate::types::builder_delegation::DashboardLink;

#[derive(Debug, Clone)]
pub struct DashboardLinkBuilder {
    internal: DashboardLink,
    errors: Vec<cog::BuildError>,
}

impl DashboardLinkBuilder {
    pub fn new() -> Self {
        Self {
            internal: DashboardLink::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for DashboardLinkBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl DashboardLinkBuilder {
    pub fn title(mut self, title: String) -> Self {
        self.internal.title = title;

        self
    }
}

impl DashboardLinkBuilder {
    pub fn url(mut self, url: String) -> Self {
        self.internal.url = url;

        self
    }
}

impl cog::Builder<DashboardLink> for DashboardLinkBuilder {
    fn build(&self) -> Result<DashboardLink, Vec<cog::BuildError>> {
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
    pub fn id(mut self, id: i64) -> Self {
        self.internal.id = id;

        self
    }
}

impl DashboardBuilder {
    pub fn title(mut self, title: String) -> Self {
        self.internal.title = title;

        self
    }
}

/// will be expanded to []cog.Builder<DashboardLink>
impl DashboardBuilder {
    pub fn links(mut self, links: Vec<impl cog::Builder<DashboardLink>>) -> Self {
        let mut built0 = Vec::new();
        for item0 in links {
            let built1 = match item0.build() {
                Ok(val) => val,
                Err(mut err) => {
                    self.errors.append(&mut err);
                    return self;
                }
            };
            built0.push(built1);
        }
        self.internal.links = built0;

        self
    }
}

/// will be expanded to [][]cog.Builder<DashboardLink>
impl DashboardBuilder {
    pub fn links_of_links(
        mut self,
        links_of_links: Vec<Vec<impl cog::Builder<DashboardLink>>>,
    ) -> Self {
        let mut built0 = Vec::new();
        for item0 in links_of_links {
            let mut built1 = Vec::new();
            for item1 in item0 {
                let built2 = match item1.build() {
                    Ok(val) => val,
                    Err(mut err) => {
                        self.errors.append(&mut err);
                        return self;
                    }
                };
                built1.push(built2);
            }
            built0.push(built1);
        }
        self.internal.links_of_links = built0;

        self
    }
}

/// will be expanded to cog.Builder<DashboardLink>
impl DashboardBuilder {
    pub fn single_link(mut self, single_link: impl cog::Builder<DashboardLink>) -> Self {
        let built0 = match single_link.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.single_link = built0;

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
