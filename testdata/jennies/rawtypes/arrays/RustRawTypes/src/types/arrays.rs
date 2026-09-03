use serde::{Deserialize, Serialize};

/// List of tags, maybe?
pub type ArrayOfStrings = Vec<String>;

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct SomeStruct {
    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,
}

pub type ArrayOfRefs = Vec<SomeStruct>;

pub type ArrayOfArrayOfNumbers = Vec<Vec<i64>>;
