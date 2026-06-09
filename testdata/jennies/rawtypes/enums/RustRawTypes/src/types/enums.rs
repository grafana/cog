use serde::{Deserialize, Serialize};
use serde_repr::Deserialize_repr;
use serde_repr::Serialize_repr;

/// This is a very interesting string enum.
#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Eq, Hash, Default)]
pub enum Operator {
    #[default]
    #[serde(rename = ">")]
    GreaterThan,
    #[serde(rename = "<")]
    LessThan,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Eq, Hash, Default)]
pub enum TableSortOrder {
    #[default]
    #[serde(rename = "asc")]
    Asc,
    #[serde(rename = "desc")]
    Desc,
}

#[derive(Serialize, Deserialize, Debug, Clone, PartialEq, Eq, Hash, Default)]
pub enum LogsSortOrder {
    #[default]
    #[serde(rename = "time_asc")]
    Asc,
    #[serde(rename = "time_desc")]
    Desc,
}

/// 0 for no shared crosshair or tooltip (default).
/// 1 for shared crosshair.
/// 2 for shared crosshair AND shared tooltip.
#[derive(Serialize_repr, Deserialize_repr, Debug, Clone, Copy, PartialEq, Eq, Hash, Default)]
#[repr(i8)]
pub enum DashboardCursorSync {
    #[default]
    Off = 0,
    Crosshair = 1,
    Tooltip = 2,
}
