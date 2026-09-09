use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct Options {
    pub timeseries_option: String,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct FieldConfig {
    pub timeseries_field_config_option: String,
}
