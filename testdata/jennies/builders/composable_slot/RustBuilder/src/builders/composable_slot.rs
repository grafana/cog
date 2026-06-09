use crate::cog;
use crate::cog::variants;
use crate::types::composable_slot::Dashboard;

#[derive(Debug, Clone)]
pub struct LokiBuilderBuilder {
    internal: Dashboard,
    errors: Vec<cog::BuildError>,
}

impl LokiBuilderBuilder {
    pub fn new() -> Self {
        Self {
            internal: Dashboard::default(),
            errors: Vec::new(),
        }
    }
}

impl Default for LokiBuilderBuilder {
    fn default() -> Self {
        Self::new()
    }
}

impl LokiBuilderBuilder {
    pub fn target(mut self, target: impl cog::Builder<Box<dyn variants::Dataquery>>) -> Self {
        let built0 = match target.build() {
            Ok(val) => val,
            Err(mut err) => {
                self.errors.append(&mut err);
                return self;
            }
        };
        self.internal.target = built0;

        self
    }
}

impl LokiBuilderBuilder {
    pub fn targets(
        mut self,
        targets: Vec<impl cog::Builder<Box<dyn variants::Dataquery>>>,
    ) -> Self {
        let mut built0 = Vec::new();
        for item0 in targets {
            let built1 = match item0.build() {
                Ok(val) => val,
                Err(mut err) => {
                    self.errors.append(&mut err);
                    return self;
                }
            };
            built0.push(built1);
        }
        self.internal.targets = built0;

        self
    }
}

impl cog::Builder<Dashboard> for LokiBuilderBuilder {
    fn build(&self) -> Result<Dashboard, Vec<cog::BuildError>> {
        if !self.errors.is_empty() {
            return Err(self.errors.clone());
        }

        Ok(self.internal.clone())
    }
}
