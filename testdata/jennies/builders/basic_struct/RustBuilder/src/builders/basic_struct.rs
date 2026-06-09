use crate::cog;
use crate::types::basic_struct::SomeStruct;

/// SomeStruct, to hold data.
#[derive(Debug, Clone)]
pub struct SomeStructBuilder {
    internal: SomeStruct,
}

impl SomeStructBuilder {
    pub fn new() -> Self {
        Self {
            internal: SomeStruct::default(),
        }
    }
}

impl Default for SomeStructBuilder {
    fn default() -> Self {
        Self::new()
    }
}

/// id identifies something. Weird, right?
impl SomeStructBuilder {
    pub fn id(mut self, id: i64) -> Self {
        self.internal.id = id;

        self
    }
}

impl SomeStructBuilder {
    pub fn uid(mut self, uid: String) -> Self {
        self.internal.uid = uid;

        self
    }
}

impl SomeStructBuilder {
    pub fn tags(mut self, tags: Vec<String>) -> Self {
        self.internal.tags = tags;

        self
    }
}

/// This thing could be live.
/// Or maybe not.
impl SomeStructBuilder {
    pub fn live_now(mut self, live_now: bool) -> Self {
        self.internal.live_now = live_now;

        self
    }
}

impl cog::Builder<SomeStruct> for SomeStructBuilder {
    fn build(&self) -> Result<SomeStruct, Vec<cog::BuildError>> {
        Ok(self.internal.clone())
    }
}
