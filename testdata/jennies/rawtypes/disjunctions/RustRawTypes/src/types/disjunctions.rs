use serde::{Deserialize, Serialize};

/// Refresh rate or disabled.
#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum RefreshRate {
    String(String),
    Bool(bool),
}

pub type StringOrNull = Option<String>;

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct SomeStruct {
    #[serde(rename = "Type")]
    #[serde(default = "default_some_struct_type")]
    pub r#type: String,

    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,
}

impl Default for SomeStruct {
    fn default() -> Self {
        Self {
            r#type: "some-struct".to_string(),
            field_any: Default::default(),
        }
    }
}

fn default_some_struct_type() -> String {
    "some-struct".to_string()
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum BoolOrRef {
    Bool(bool),
    SomeStruct(SomeStruct),
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct SomeOtherStruct {
    #[serde(rename = "Type")]
    #[serde(default = "default_some_other_struct_type")]
    pub r#type: String,

    #[serde(rename = "Foo")]
    pub foo: Vec<u8>,
}

impl Default for SomeOtherStruct {
    fn default() -> Self {
        Self {
            r#type: "some-other-struct".to_string(),
            foo: Default::default(),
        }
    }
}

fn default_some_other_struct_type() -> String {
    "some-other-struct".to_string()
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct YetAnotherStruct {
    #[serde(rename = "Type")]
    #[serde(default = "default_yet_another_struct_type")]
    pub r#type: String,

    #[serde(rename = "Bar")]
    pub bar: u8,
}

impl Default for YetAnotherStruct {
    fn default() -> Self {
        Self {
            r#type: "yet-another-struct".to_string(),
            bar: Default::default(),
        }
    }
}

fn default_yet_another_struct_type() -> String {
    "yet-another-struct".to_string()
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum SeveralRefs {
    SomeStruct(SomeStruct),
    SomeOtherStruct(SomeOtherStruct),
    YetAnotherStruct(YetAnotherStruct),
}
