use crate::types::sandbox;

/// some_struct_converter accepts a `SomeStruct` object and generates the Rust code to build this object using builders.
pub fn some_struct_converter(input: &sandbox::SomeStruct) -> String {
    let mut calls: Vec<String> = vec!["sandbox::SomeStructBuilder::new()".to_string()];
    if input.editable {
        let mut buffer = String::new();
        buffer.push_str("editable(");
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.editable {
        let mut buffer = String::new();
        buffer.push_str("readonly(");
        buffer.push(')');
        calls.push(buffer);
    }
    if input.auto_refresh.is_some() && input.auto_refresh == Some(true) {
        let mut buffer = String::new();
        buffer.push_str("auto_refresh(");
        buffer.push(')');
        calls.push(buffer);
    }
    if input.auto_refresh.is_some() && input.auto_refresh == Some(false) {
        let mut buffer = String::new();
        buffer.push_str("no_auto_refresh(");
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
