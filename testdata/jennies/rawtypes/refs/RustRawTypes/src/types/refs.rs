use crate::types::otherpkg::SomeDistantStruct;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct SomeStruct {
    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,
}

pub type RefToSomeStruct = SomeStruct;

pub type RefToSomeStructFromOtherPackage = SomeDistantStruct;
