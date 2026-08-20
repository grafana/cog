use crate::types::nullable_map_assignment;

/// some_struct_converter accepts a `SomeStruct` object and generates the Rust code to build this object using builders.
pub fn some_struct_converter(input: &nullable_map_assignment::SomeStruct) -> String {
    let mut calls: Vec<String> =
        vec!["nullable_map_assignment::SomeStructBuilder::new()".to_string()];
    if !input.config.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("config(");
        let mut arg0 = String::from("std::collections::HashMap::from([");
        for (key, arg1) in &input.config {
            let tmp_config_arg1 = format!("{:?}.to_string()", arg1);
            let tmp_config_arg1_key = format!("{:?}.to_string()", key);
            arg0.push_str(&format!("({}, {}), ", tmp_config_arg1_key, tmp_config_arg1));
        }
        arg0.push_str("])");
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
