use serde::{Deserialize, Serialize};
use std::collections::HashMap;

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct SomeStruct {
    pub id: u64,

    #[serde(rename = "maybeId")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub maybe_id: Option<u64>,

    #[serde(rename = "greaterThanZero")]
    pub greater_than_zero: u64,

    pub negative: i64,

    pub title: String,

    pub labels: HashMap<String, String>,

    pub tags: Vec<String>,

    pub regex: String,

    #[serde(rename = "negativeRegex")]
    pub negative_regex: String,

    #[serde(rename = "minMaxList")]
    pub min_max_list: Vec<String>,

    #[serde(rename = "uniqueList")]
    pub unique_list: Vec<String>,

    #[serde(rename = "fullConstraintList")]
    pub full_constraint_list: Vec<i64>,
}
