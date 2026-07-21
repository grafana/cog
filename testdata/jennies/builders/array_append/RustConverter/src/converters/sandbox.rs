use crate::types::sandbox;

/// some_struct_converter accepts a `SomeStruct` object and generates the Rust code to build this object using builders.
pub fn some_struct_converter(input: &sandbox::SomeStruct) -> String {
    let mut calls: Vec<String> = vec!["sandbox::SomeStructBuilder::new()".to_string()];
    if !input.tags.is_empty() {
        for item in &input.tags {
            let mut buffer = String::new();
            buffer.push_str("tags(");
            let arg0 = format!("{:?}.to_string()", item);
            buffer.push_str(&arg0);
            buffer.push(')');
            calls.push(buffer);
        }
    }

    calls.join("\n    .")
}
