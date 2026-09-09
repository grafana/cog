use crate::types::constructor_initializations;

/// some_panel_converter accepts a `SomePanel` object and generates the Rust code to build this object using builders.
pub fn some_panel_converter(input: &constructor_initializations::SomePanel) -> String {
    let mut calls: Vec<String> =
        vec!["constructor_initializations::SomePanelBuilder::new()".to_string()];
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
