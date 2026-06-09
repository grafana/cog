use serde::{Deserialize, Serialize};

/// This
/// is
/// a
/// comment
#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct SomeStruct {
    /// Anything can go in there.
    /// Really, anything.
    #[serde(rename = "FieldAny")]
    pub field_any: serde_json::Value,

    #[serde(rename = "FieldBool")]
    pub field_bool: bool,

    #[serde(rename = "FieldBytes")]
    pub field_bytes: Vec<u8>,

    #[serde(rename = "FieldString")]
    pub field_string: String,

    #[serde(rename = "FieldStringWithConstantValue")]
    pub field_string_with_constant_value: String,

    #[serde(rename = "FieldFloat32")]
    pub field_float32: f32,

    #[serde(rename = "FieldFloat64")]
    pub field_float64: f64,

    #[serde(rename = "FieldUint8")]
    pub field_uint8: u8,

    #[serde(rename = "FieldUint16")]
    pub field_uint16: u16,

    #[serde(rename = "FieldUint32")]
    pub field_uint32: u32,

    #[serde(rename = "FieldUint64")]
    pub field_uint64: u64,

    #[serde(rename = "FieldInt8")]
    pub field_int8: i8,

    #[serde(rename = "FieldInt16")]
    pub field_int16: i16,

    #[serde(rename = "FieldInt32")]
    pub field_int32: i32,

    #[serde(rename = "FieldInt64")]
    pub field_int64: i64,
}

impl Default for SomeStruct {
    fn default() -> Self {
        Self {
            field_any: Default::default(),
            field_bool: Default::default(),
            field_bytes: Default::default(),
            field_string: Default::default(),
            field_string_with_constant_value: "auto".to_string(),
            field_float32: Default::default(),
            field_float64: Default::default(),
            field_uint8: Default::default(),
            field_uint16: Default::default(),
            field_uint32: Default::default(),
            field_uint64: Default::default(),
            field_int8: Default::default(),
            field_int16: Default::default(),
            field_int32: Default::default(),
            field_int64: Default::default(),
        }
    }
}
