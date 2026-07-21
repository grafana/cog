use crate::types::known_any;

/// some_struct_converter accepts a `SomeStruct` object and generates the Rust code to build this object using builders.
pub fn some_struct_converter(input: &known_any::SomeStruct) -> String {
    let mut calls: Vec<String> = vec!["known_any::SomeStructBuilder::new()".to_string()];
    if input.config.is_some()
        && !input.config.as_ref().unwrap()["title"]
            .as_str()
            .unwrap_or_default()
            .is_empty()
    {
        let mut buffer = String::new();
        buffer.push_str("title(");
        let arg0 = format!(
            "{:?}.to_string()",
            input.config.as_ref().unwrap()["title"]
                .as_str()
                .unwrap_or_default()
        );
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
