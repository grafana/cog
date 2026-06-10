use crate::types::sandbox;

/// some_struct_converter accepts a `SomeStruct` object and generates the Rust code to build this object using builders.
pub fn some_struct_converter(input: &sandbox::SomeStruct) -> String {
    let mut calls: Vec<String> = vec!["sandbox::SomeStructBuilder::new()".to_string()];
    if !input.annotations.is_empty() {
        for (key, value) in &input.annotations {
            let mut buffer = String::new();
            buffer.push_str("annotations(");
            let arg0 = format!("{:?}.to_string()", key);
            buffer.push_str(&arg0);
            buffer.push_str(", ");
            let arg1 = format!("{:?}.to_string()", value);
            buffer.push_str(&arg1);
            buffer.push(')');
            calls.push(buffer);
        }
    }

    calls.join("\n    .")
}
