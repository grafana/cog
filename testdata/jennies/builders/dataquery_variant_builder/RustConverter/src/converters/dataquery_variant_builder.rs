use crate::types::dataquery_variant_builder;

/// loki_builder_converter accepts a `LokiBuilder` object and generates the Rust code to build this object using builders.
pub fn loki_builder_converter(input: &dataquery_variant_builder::Loki) -> String {
    let mut calls: Vec<String> =
        vec!["dataquery_variant_builder::LokiBuilderBuilder::new()".to_string()];
    if !input.expr.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("expr(");
        let arg0 = format!("{:?}.to_string()", input.expr);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
