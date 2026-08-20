use crate::cog;
use crate::types::discriminator_without_option;

/// no_show_field_option_converter accepts a `NoShowFieldOption` object and generates the Rust code to build this object using builders.
pub fn no_show_field_option_converter(
    input: &discriminator_without_option::NoShowFieldOption,
) -> String {
    let mut calls: Vec<String> =
        vec!["discriminator_without_option::NoShowFieldOptionBuilder::new()".to_string()];
    if !input.text.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("text(");
        let arg0 = format!("{:?}.to_string()", input.text);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}

/// show_field_option_converter accepts a `ShowFieldOption` object and generates the Rust code to build this object using builders.
pub fn show_field_option_converter(
    input: &discriminator_without_option::ShowFieldOption,
) -> String {
    let mut calls: Vec<String> =
        vec!["discriminator_without_option::ShowFieldOptionBuilder::new()".to_string()];
    {
        let mut buffer = String::new();
        buffer.push_str("field(");
        let arg0 = cog::dump(&input.field);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.text.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("text(");
        let arg0 = format!("{:?}.to_string()", input.text);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
