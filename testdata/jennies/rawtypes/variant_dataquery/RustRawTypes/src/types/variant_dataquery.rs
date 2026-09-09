use crate::cog::variants;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct Query {
    pub expr: String,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub instant: Option<bool>,
}

impl variants::Dataquery for Query {
    fn dataquery_type(&self) -> String {
        "prometheus".to_string()
    }

    fn dataquery_equals(&self, other: &dyn variants::Dataquery) -> bool {
        match (
            variants::DataquerySerialize::to_json_value(self),
            other.to_json_value(),
        ) {
            (Ok(a), Ok(b)) => a == b,
            _ => false,
        }
    }
}
