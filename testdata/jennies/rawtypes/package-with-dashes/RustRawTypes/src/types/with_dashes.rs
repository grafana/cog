use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct SomeStruct {
    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,
}

/// Refresh rate or disabled.
#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum RefreshRate {
    String(String),
    Bool(bool),
}
