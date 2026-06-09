use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum DisjunctionOfScalarsAndRefs {
    String(String),
    Bool(bool),
    VecString(Vec<String>),
    MyRefA(MyRefA),
    MyRefB(MyRefB),
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct MyRefA {
    pub foo: String,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct MyRefB {
    pub bar: i64,
}
