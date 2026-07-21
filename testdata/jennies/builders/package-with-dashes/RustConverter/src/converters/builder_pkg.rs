use crate::types::with_dashes;

/// some_nice_builder_converter accepts a `SomeNiceBuilder` object and generates the Rust code to build this object using builders.
pub fn some_nice_builder_converter(input: &with_dashes::SomeStruct) -> String {
    let mut calls: Vec<String> = vec!["builder_pkg::SomeNiceBuilderBuilder::new()".to_string()];
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
