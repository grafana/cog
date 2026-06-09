use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct SomeStruct {
    #[serde(rename = "fieldBool")]
    pub field_bool: bool,

    #[serde(rename = "fieldString")]
    pub field_string: String,

    #[serde(rename = "FieldStringWithConstantValue")]
    #[serde(default = "default_some_struct_field_string_with_constant_value")]
    pub field_string_with_constant_value: String,

    #[serde(rename = "FieldFloat32")]
    pub field_float32: f32,

    #[serde(rename = "FieldInt32")]
    pub field_int32: i32,
}

impl Default for SomeStruct {
    fn default() -> Self {
        Self {
            field_bool: true,
            field_string: "foo".to_string(),
            field_string_with_constant_value: "auto".to_string(),
            field_float32: 42.42,
            field_int32: 42,
        }
    }
}

fn default_some_struct_field_string_with_constant_value() -> String {
    "auto".to_string()
}
