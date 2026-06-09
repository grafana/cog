use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum DisjunctionWithoutDiscriminator {
    TypeA(TypeA),
    TypeB(TypeB),
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct TypeA {
    #[serde(rename = "fieldA")]
    pub field_a: String,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct TypeB {
    #[serde(rename = "fieldB")]
    pub field_b: i64,
}
