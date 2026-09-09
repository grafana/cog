use crate::types::basic_struct;

/// some_struct_converter accepts a `SomeStruct` object and generates the Rust code to build this object using builders.
pub fn some_struct_converter(input: &basic_struct::SomeStruct) -> String {
    let mut calls: Vec<String> = vec!["basic_struct::SomeStructBuilder::new()".to_string()];
    {
        let mut buffer = String::new();
        buffer.push_str("id(");
        let arg0 = format!("{:?}", input.id);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.uid.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("uid(");
        let arg0 = format!("{:?}.to_string()", input.uid);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.tags.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("tags(");
        let mut tmp_arg0: Vec<String> = Vec::new();
        for arg1 in &input.tags {
            let tmp_tags_arg1 = format!("{:?}.to_string()", arg1);
            tmp_arg0.push(tmp_tags_arg1);
        }
        let arg0 = format!("vec![{}]", tmp_arg0.join(", "));
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    {
        let mut buffer = String::new();
        buffer.push_str("live_now(");
        let arg0 = format!("{:?}", input.live_now);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
