use crate::types::initialization_safeguards;

/// some_panel_converter accepts a `SomePanel` object and generates the Rust code to build this object using builders.
pub fn some_panel_converter(input: &initialization_safeguards::SomePanel) -> String {
    let mut calls: Vec<String> =
        vec!["initialization_safeguards::SomePanelBuilder::new()".to_string()];
    if !input.title.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("title(");
        let arg0 = format!("{:?}.to_string()", input.title);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if input.options.is_some() {
        let mut buffer = String::new();
        buffer.push_str("show_legend(");
        let arg0 = format!("{:?}", input.options.as_ref().unwrap().legend.show);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
