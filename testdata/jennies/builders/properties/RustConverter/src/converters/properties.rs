use crate::types::properties;

/// some_struct_converter accepts a `SomeStruct` object and generates the Rust code to build this object using builders.
pub fn some_struct_converter(input: &properties::SomeStruct) -> String {
    let mut calls: Vec<String> = vec!["properties::SomeStructBuilder::new()".to_string()];
    {
        let mut buffer = String::new();
        buffer.push_str("id(");
        let arg0 = format!("{:?}", input.id);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
