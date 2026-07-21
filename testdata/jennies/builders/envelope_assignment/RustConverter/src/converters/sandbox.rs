use crate::types::sandbox;

/// dashboard_converter accepts a `Dashboard` object and generates the Rust code to build this object using builders.
pub fn dashboard_converter(input: &sandbox::Dashboard) -> String {
    let mut calls: Vec<String> = vec!["sandbox::DashboardBuilder::new()".to_string()];
    if !input.variables.is_empty() {
        for item in &input.variables {
            let mut buffer = String::new();
            buffer.push_str("with_variable(");
            let arg0 = format!("{:?}.to_string()", item.name);
            buffer.push_str(&arg0);
            buffer.push_str(", ");
            let arg1 = format!("{:?}.to_string()", item.value);
            buffer.push_str(&arg1);
            buffer.push(')');
            calls.push(buffer);
        }
    }

    calls.join("\n    .")
}
