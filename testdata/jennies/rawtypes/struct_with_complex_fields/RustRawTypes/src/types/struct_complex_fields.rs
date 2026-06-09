use serde::{Deserialize, Serialize};
use std::collections::HashMap;

/// This struct does things.
#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct SomeStruct {
    #[serde(rename = "FieldRef")]
    pub field_ref: SomeOtherStruct,

    #[serde(rename = "FieldDisjunctionOfScalars")]
    pub field_disjunction_of_scalars: SomeStructFieldDisjunctionOfScalars,

    #[serde(rename = "FieldMixedDisjunction")]
    pub field_mixed_disjunction: SomeStructFieldMixedDisjunction,

    #[serde(rename = "FieldDisjunctionWithNull")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub field_disjunction_with_null: Option<String>,

    #[serde(rename = "Operator")]
    pub operator: SomeStructOperator,

    #[serde(rename = "FieldArrayOfStrings")]
    pub field_array_of_strings: Vec<String>,

    #[serde(rename = "FieldMapOfStringToString")]
    pub field_map_of_string_to_string: HashMap<String, String>,

    #[serde(rename = "FieldAnonymousStruct")]
    pub field_anonymous_struct: StructComplexFieldsSomeStructFieldAnonymousStruct,

    #[serde(rename = "fieldRefToConstant")]
    pub field_ref_to_constant: String,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum SomeStructFieldDisjunctionOfScalars {
    String(String),
    Bool(bool),
}

impl Default for SomeStructFieldDisjunctionOfScalars {
    fn default() -> Self {
        Self::String(Default::default())
    }
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum SomeStructFieldMixedDisjunction {
    String(String),
    SomeOtherStruct(SomeOtherStruct),
}

impl Default for SomeStructFieldMixedDisjunction {
    fn default() -> Self {
        Self::String(Default::default())
    }
}

pub const CONNECTION_PATH: &str = "straight";

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct SomeOtherStruct {
    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct StructComplexFieldsSomeStructFieldAnonymousStruct {
    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Eq, Hash, Default)]
pub enum SomeStructOperator {
    #[default]
    #[serde(rename = ">")]
    GreaterThan,
    #[serde(rename = "<")]
    LessThan,
}
