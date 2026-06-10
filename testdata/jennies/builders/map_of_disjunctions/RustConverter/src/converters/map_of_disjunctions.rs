use crate::types::map_of_disjunctions;

/// element_converter accepts a `Element` object and generates the Rust code to build this object using builders.
pub fn element_converter(input: &map_of_disjunctions::Element) -> String {
    let mut calls: Vec<String> = vec!["map_of_disjunctions::ElementBuilder::new()".to_string()];
    if input.panel.is_some() {
        let mut buffer = String::new();
        buffer.push_str("panel(");
        let arg0 = panel_converter(input.panel.as_ref().unwrap());
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if input.library_panel.is_some() {
        let mut buffer = String::new();
        buffer.push_str("library_panel(");
        let arg0 = library_panel_converter(input.library_panel.as_ref().unwrap());
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}

/// panel_converter accepts a `Panel` object and generates the Rust code to build this object using builders.
pub fn panel_converter(input: &map_of_disjunctions::Panel) -> String {
    let mut calls: Vec<String> = vec!["map_of_disjunctions::PanelBuilder::new()".to_string()];
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

/// library_panel_converter accepts a `LibraryPanel` object and generates the Rust code to build this object using builders.
pub fn library_panel_converter(input: &map_of_disjunctions::LibraryPanel) -> String {
    let mut calls: Vec<String> =
        vec!["map_of_disjunctions::LibraryPanelBuilder::new()".to_string()];
    if !input.text.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("text(");
        let arg0 = format!("{:?}.to_string()", input.text);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}

/// dashboard_converter accepts a `Dashboard` object and generates the Rust code to build this object using builders.
pub fn dashboard_converter(input: &map_of_disjunctions::Dashboard) -> String {
    let mut calls: Vec<String> = vec!["map_of_disjunctions::DashboardBuilder::new()".to_string()];
    if !input.panels.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("panels(");
        let mut arg0 = String::from("std::collections::HashMap::from([");
        for (key, arg1) in &input.panels {
            let tmp_panels_arg1 = element_converter(arg1);
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

/// panel_or_library_panel_converter accepts a `PanelOrLibraryPanel` object and generates the Rust code to build this object using builders.
pub fn panel_or_library_panel_converter(
    input: &map_of_disjunctions::PanelOrLibraryPanel,
) -> String {
    let mut calls: Vec<String> =
        vec!["map_of_disjunctions::PanelOrLibraryPanelBuilder::new()".to_string()];
    if input.panel.is_some() {
        let mut buffer = String::new();
        buffer.push_str("panel(");
        let arg0 = panel_converter(input.panel.as_ref().unwrap());
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if input.library_panel.is_some() {
        let mut buffer = String::new();
        buffer.push_str("library_panel(");
        let arg0 = library_panel_converter(input.library_panel.as_ref().unwrap());
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
