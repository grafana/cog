use crate::cog::variants;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct Dashboard {
    pub title: String,

    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub panels: Vec<Panel>,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct DataSourceRef {
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub r#type: Option<String>,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub uid: Option<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct FieldConfigSource {
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub defaults: Option<FieldConfig>,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct FieldConfig {
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub unit: Option<String>,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub custom: Option<serde_json::Value>,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct Panel {
    pub title: String,

    pub r#type: String,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub datasource: Option<DataSourceRef>,

    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub options: Option<serde_json::Value>,

    #[serde(deserialize_with = "crate::cog::variants::deserialize_dataquery_vec")]
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub targets: Vec<Box<dyn variants::Dataquery>>,

    #[serde(rename = "fieldConfig")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub field_config: Option<FieldConfigSource>,
}
