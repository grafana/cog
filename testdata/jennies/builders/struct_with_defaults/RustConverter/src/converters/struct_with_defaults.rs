use crate::cog;
use crate::types::struct_with_defaults;

/// nested_struct_converter accepts a `NestedStruct` object and generates the Rust code to build this object using builders.
pub fn nested_struct_converter(input: &struct_with_defaults::NestedStruct) -> String {
    let mut calls: Vec<String> =
        vec!["struct_with_defaults::NestedStructBuilder::new()".to_string()];
    if !input.string_val.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("string_val(");
        let arg0 = format!("{:?}.to_string()", input.string_val);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    {
        let mut buffer = String::new();
        buffer.push_str("int_val(");
        let arg0 = format!("{:?}", input.int_val);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}

/// struct_converter accepts a `Struct` object and generates the Rust code to build this object using builders.
pub fn struct_converter(input: &struct_with_defaults::Struct) -> String {
    let mut calls: Vec<String> = vec!["struct_with_defaults::StructBuilder::new()".to_string()];
    {
        let mut buffer = String::new();
        buffer.push_str("all_fields(");
        let arg0 = nested_struct_converter(&input.all_fields);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    {
        let mut buffer = String::new();
        buffer.push_str("partial_fields(");
        let arg0 = nested_struct_converter(&input.partial_fields);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    {
        let mut buffer = String::new();
        buffer.push_str("empty_fields(");
        let arg0 = nested_struct_converter(&input.empty_fields);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    {
        let mut buffer = String::new();
        buffer.push_str("complex_field(");
        let arg0 = cog::dump(&input.complex_field);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    {
        let mut buffer = String::new();
        buffer.push_str("partial_complex_field(");
        let arg0 = cog::dump(&input.partial_complex_field);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
