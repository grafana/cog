use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Eq, Hash, Default)]
pub enum Enum {
    #[default]
    #[serde(rename = "ValueA")]
    EnumValueA,
    #[serde(rename = "ValueB")]
    EnumValueB,
    #[serde(rename = "ValueC")]
    EnumValueC,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct ParentStruct {
    #[serde(rename = "myEnum")]
    pub my_enum: Enum,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct Struct {
    #[serde(rename = "myValue")]
    pub my_value: String,

    #[serde(rename = "myEnum")]
    pub my_enum: Enum,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct StructA {
    #[serde(rename = "myEnum")]
    #[serde(default = "default_struct_a_my_enum")]
    pub my_enum: Enum,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub other: Option<Enum>,
}

impl Default for StructA {
    fn default() -> Self {
        Self {
            my_enum: Enum::EnumValueA,
            other: Some(Enum::EnumValueA),
        }
    }
}

fn default_struct_a_my_enum() -> Enum {
    Enum::EnumValueA
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct StructB {
    #[serde(rename = "myEnum")]
    #[serde(default = "default_struct_b_my_enum")]
    pub my_enum: Enum,

    #[serde(rename = "myValue")]
    pub my_value: String,
}

impl Default for StructB {
    fn default() -> Self {
        Self {
            my_enum: Enum::EnumValueB,
            my_value: Default::default(),
        }
    }
}

fn default_struct_b_my_enum() -> Enum {
    Enum::EnumValueB
}
