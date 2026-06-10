use crate::cog;
use crate::types::promql;

/// func_call_expr_converter accepts a `FuncCallExpr` object and generates the Rust code to build this object using builders.
pub fn func_call_expr_converter(input: &promql::FuncCallExpr) -> String {
    let mut calls: Vec<String> = vec!["promql::FuncCallExprBuilder::new()".to_string()];
    if !input.function.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("function(");
        let arg0 = format!("{:?}.to_string()", input.function);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.args.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("args(");
        let mut tmp_arg0: Vec<String> = Vec::new();
        for arg1 in &input.args {
            let tmp_args_arg1 = cog::dump(arg1);
            tmp_arg0.push(tmp_args_arg1);
        }
        let arg0 = format!("vec![{}]", tmp_arg0.join(", "));
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
