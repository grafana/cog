use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct SomeStruct {
    #[serde(default, skip_serializing_if = "HashMap::is_empty")]
    pub options: HashMap<String, serde_json::Value>,

    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub items: Vec<String>,

    pub extra: serde_json::Value,
}

impl Default for SomeStruct {
    fn default() -> Self {
        Self {
            options: HashMap::new(),
            items: Vec::new(),
            extra: serde_json::Value::Object(serde_json::Map::new()),
        }
    }
}
