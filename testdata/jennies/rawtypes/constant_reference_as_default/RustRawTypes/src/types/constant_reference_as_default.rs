use serde::{Deserialize, Serialize};

pub const CONSTANT_REF_STRING: &str = "AString";

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct MyStruct {
    #[serde(rename = "aString")]
    #[serde(default = "default_my_struct_a_string")]
    pub a_string: String,

    #[serde(rename = "optString")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub opt_string: Option<String>,
}

impl Default for MyStruct {
    fn default() -> Self {
        Self {
            a_string: "AString".to_string(),
            opt_string: Some("AString".to_string()),
        }
    }
}

fn default_my_struct_a_string() -> String {
    "AString".to_string()
}
