use crate::cog;
use crate::types::composable_slot;

/// loki_builder_converter accepts a `LokiBuilder` object and generates the Rust code to build this object using builders.
pub fn loki_builder_converter(input: &composable_slot::Dashboard) -> String {
    let mut calls: Vec<String> = vec!["composable_slot::LokiBuilderBuilder::new()".to_string()];
    {
        let mut buffer = String::new();
        buffer.push_str("target(");
        let arg0 = cog::convert_dataquery_to_code(input.target.as_ref());
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.targets.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("targets(");
        let mut tmp_arg0: Vec<String> = Vec::new();
        for arg1 in &input.targets {
            let tmp_targets_arg1 = cog::convert_dataquery_to_code(arg1.as_ref());
            tmp_arg0.push(tmp_targets_arg1);
        }
        let arg0 = format!("vec![{}]", tmp_arg0.join(", "));
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
