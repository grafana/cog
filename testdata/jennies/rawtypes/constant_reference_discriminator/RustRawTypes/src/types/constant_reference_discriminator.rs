use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum LayoutWithValue {
    GridLayoutUsingValue(GridLayoutUsingValue),
    RowsLayoutUsingValue(RowsLayoutUsingValue),
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct GridLayoutUsingValue {
    #[serde(default = "default_grid_layout_using_value_kind")]
    pub kind: String,

    #[serde(rename = "gridLayoutProperty")]
    pub grid_layout_property: String,
}

impl Default for GridLayoutUsingValue {
    fn default() -> Self {
        Self {
            kind: "GridLayout".to_string(),
            grid_layout_property: Default::default(),
        }
    }
}

fn default_grid_layout_using_value_kind() -> String {
    "GridLayout".to_string()
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct RowsLayoutUsingValue {
    #[serde(default = "default_rows_layout_using_value_kind")]
    pub kind: String,

    #[serde(rename = "rowsLayoutProperty")]
    pub rows_layout_property: String,
}

impl Default for RowsLayoutUsingValue {
    fn default() -> Self {
        Self {
            kind: "RowsLayout".to_string(),
            rows_layout_property: Default::default(),
        }
    }
}

fn default_rows_layout_using_value_kind() -> String {
    "RowsLayout".to_string()
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
#[serde(untagged)]
pub enum LayoutWithoutValue {
    GridLayoutWithoutValue(GridLayoutWithoutValue),
    RowsLayoutWithoutValue(RowsLayoutWithoutValue),
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct GridLayoutWithoutValue {
    #[serde(default = "default_grid_layout_without_value_kind")]
    pub kind: String,

    #[serde(rename = "gridLayoutProperty")]
    pub grid_layout_property: String,
}

impl Default for GridLayoutWithoutValue {
    fn default() -> Self {
        Self {
            kind: "GridLayout".to_string(),
            grid_layout_property: Default::default(),
        }
    }
}

fn default_grid_layout_without_value_kind() -> String {
    "GridLayout".to_string()
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq)]
pub struct RowsLayoutWithoutValue {
    #[serde(default = "default_rows_layout_without_value_kind")]
    pub kind: String,

    #[serde(rename = "rowsLayoutProperty")]
    pub rows_layout_property: String,
}

impl Default for RowsLayoutWithoutValue {
    fn default() -> Self {
        Self {
            kind: "RowsLayout".to_string(),
            rows_layout_property: Default::default(),
        }
    }
}

fn default_rows_layout_without_value_kind() -> String {
    "RowsLayout".to_string()
}

pub const GRID_LAYOUT_KIND_TYPE: &str = "GridLayout";

pub const ROWS_LAYOUT_KIND_TYPE: &str = "RowsLayout";
