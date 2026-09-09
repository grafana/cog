use crate::cog;
use crate::types::some_pkg;

/// person_converter accepts a `Person` object and generates the Rust code to build this object using builders.
pub fn person_converter(input: &some_pkg::Person) -> String {
    let mut calls: Vec<String> = vec!["some_pkg::PersonBuilder::new()".to_string()];
    {
        let mut buffer = String::new();
        buffer.push_str("name(");
        let arg0 = cog::dump(&input.name);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
