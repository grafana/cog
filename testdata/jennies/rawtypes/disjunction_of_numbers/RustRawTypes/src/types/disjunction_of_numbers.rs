use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum Numbers {
    I64(i64),
    F64(f64),
    F32(f32),
}

impl Default for Numbers {
    fn default() -> Self {
        Self::I64(Default::default())
    }
}
