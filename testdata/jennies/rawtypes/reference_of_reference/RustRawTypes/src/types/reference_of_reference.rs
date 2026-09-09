use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct MyStruct {
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub field: Option<OtherStruct>,
}

pub type OtherStruct = AnotherStruct;

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct AnotherStruct {
    pub a: String,
}
