use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct Struct {
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub a: Option<MyObject>,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub b: Option<MyObject>,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub c: Option<String>,

    pub d: Vec<String>,

    pub e: HashMap<String, Option<String>>,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub f: Option<NullableFieldsStructF>,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub g: Option<String>,
}

impl Default for Struct {
    fn default() -> Self {
        Self {
            a: Default::default(),
            b: Default::default(),
            c: Default::default(),
            d: Default::default(),
            e: Default::default(),
            f: Default::default(),
            g: Some("hey".to_string()),
        }
    }
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct MyObject {
    pub field: String,
}

pub const CONSTANT_REF: &str = "hey";

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct NullableFieldsStructF {
    pub a: String,
}
