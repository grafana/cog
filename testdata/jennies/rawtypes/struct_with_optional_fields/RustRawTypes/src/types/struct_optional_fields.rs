use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct SomeStruct {
    #[serde(rename = "FieldRef")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub field_ref: Option<SomeOtherStruct>,

    #[serde(rename = "FieldString")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub field_string: Option<String>,

    #[serde(rename = "Operator")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub operator: Option<SomeStructOperator>,

    #[serde(rename = "FieldArrayOfStrings")]
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub field_array_of_strings: Vec<String>,

    #[serde(rename = "FieldAnonymousStruct")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub field_anonymous_struct: Option<StructOptionalFieldsSomeStructFieldAnonymousStruct>,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct SomeOtherStruct {
    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct StructOptionalFieldsSomeStructFieldAnonymousStruct {
    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Eq, Hash)]
pub enum SomeStructOperator {
    #[serde(rename = ">")]
    GreaterThan,
    #[serde(rename = "<")]
    LessThan,
}
