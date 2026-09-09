use crate::types::map_of_builders;

/// panel_converter accepts a `Panel` object and generates the Rust code to build this object using builders.
pub fn panel_converter(input: &map_of_builders::Panel) -> String {
    let mut calls: Vec<String> = vec!["map_of_builders::PanelBuilder::new()".to_string()];
    if !input.title.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("title(");
        let arg0 = format!("{:?}.to_string()", input.title);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}

/// dashboard_converter accepts a `Dashboard` object and generates the Rust code to build this object using builders.
pub fn dashboard_converter(input: &map_of_builders::Dashboard) -> String {
    let mut calls: Vec<String> = vec!["map_of_builders::DashboardBuilder::new()".to_string()];
    if !input.panels.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("panels(");
        let mut arg0 = String::from("std::collections::HashMap::from([");
        for (key, arg1) in &input.panels {
            let tmp_panels_arg1 = panel_converter(arg1);
            let tmp_panels_arg1_key = format!("{:?}.to_string()", key);
            arg0.push_str(&format!("({}, {}), ", tmp_panels_arg1_key, tmp_panels_arg1));
        }
        arg0.push_str("])");
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
