use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// String to... something.
pub type MapOfStringToAny = HashMap<String, serde_json::Value>;

pub type MapOfStringToString = HashMap<String, String>;

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct SomeStruct {
    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,
}

pub type MapOfStringToRef = HashMap<String, SomeStruct>;

pub type MapOfStringToMapOfStringToBool = HashMap<String, HashMap<String, bool>>;
