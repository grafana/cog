use crate::types::sandbox;

/// some_struct_converter accepts a `SomeStruct` object and generates the Rust code to build this object using builders.
pub fn some_struct_converter(input: &sandbox::SomeStruct) -> String {
    let mut calls: Vec<String> = vec![format!(
        "sandbox::SomeStructBuilder::new({})",
        format!("{:?}.to_string()", input.title)
    )];
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
