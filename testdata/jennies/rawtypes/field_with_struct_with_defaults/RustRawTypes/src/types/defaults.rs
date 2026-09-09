use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct NestedStruct {
    #[serde(rename = "stringVal")]
    pub string_val: String,

    #[serde(rename = "intVal")]
    pub int_val: i64,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct Struct {
    #[serde(rename = "allFields")]
    pub all_fields: NestedStruct,

    #[serde(rename = "partialFields")]
    pub partial_fields: NestedStruct,

    #[serde(rename = "emptyFields")]
    pub empty_fields: NestedStruct,

    #[serde(rename = "complexField")]
    pub complex_field: DefaultsStructComplexField,

    #[serde(rename = "partialComplexField")]
    pub partial_complex_field: DefaultsStructPartialComplexField,
}

impl Default for Struct {
    fn default() -> Self {
        Self {
            all_fields: NestedStruct {
                string_val: "hello".to_string(),
                int_val: 3,
            },
            partial_fields: NestedStruct {
                int_val: 3,
                ..Default::default()
            },
            empty_fields: Default::default(),
            complex_field: DefaultsStructComplexField {
                uid: "myUID".to_string(),
                nested: DefaultsStructComplexFieldNested {
                    nested_val: "nested".to_string(),
                },
                array: vec!["hello".to_string()],
            },
            partial_complex_field: DefaultsStructPartialComplexField::default(),
        }
    }
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct DefaultsStructComplexFieldNested {
    #[serde(rename = "nestedVal")]
    pub nested_val: String,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct DefaultsStructComplexField {
    pub uid: String,

    pub nested: DefaultsStructComplexFieldNested,

    pub array: Vec<String>,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Default)]
pub struct DefaultsStructPartialComplexField {
    pub uid: String,

    #[serde(rename = "intVal")]
    pub int_val: i64,
}
