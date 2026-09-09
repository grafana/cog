use serde::{Deserialize, Serialize};

pub type ObjTime = String;

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct ObjWithTimeField {
    #[serde(rename = "registeredAt")]
    pub registered_at: String,

    pub duration: String,
}
