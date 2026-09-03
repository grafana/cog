use crate::types::panelbuilder;

/// panel_converter accepts a `Panel` object and generates the Rust code to build this object using builders.
pub fn panel_converter(input: &panelbuilder::Panel) -> String {
    let mut calls: Vec<String> = vec!["panelbuilder::PanelBuilder::new()".to_string()];
    if input.only_from_this_dashboard {
        let mut buffer = String::new();
        buffer.push_str("only_from_this_dashboard(");
        let arg0 = format!("{:?}", input.only_from_this_dashboard);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if input.only_in_time_range {
        let mut buffer = String::new();
        buffer.push_str("only_in_time_range(");
        let arg0 = format!("{:?}", input.only_in_time_range);
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
    if input.limit != 10 {
        let mut buffer = String::new();
        buffer.push_str("limit(");
        let arg0 = format!("{:?}", input.limit);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.show_user {
        let mut buffer = String::new();
        buffer.push_str("show_user(");
        let arg0 = format!("{:?}", input.show_user);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.show_time {
        let mut buffer = String::new();
        buffer.push_str("show_time(");
        let arg0 = format!("{:?}", input.show_time);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.show_tags {
        let mut buffer = String::new();
        buffer.push_str("show_tags(");
        let arg0 = format!("{:?}", input.show_tags);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.navigate_to_panel {
        let mut buffer = String::new();
        buffer.push_str("navigate_to_panel(");
        let arg0 = format!("{:?}", input.navigate_to_panel);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.navigate_before.is_empty() && input.navigate_before != "10m" {
        let mut buffer = String::new();
        buffer.push_str("navigate_before(");
        let arg0 = format!("{:?}.to_string()", input.navigate_before);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.navigate_after.is_empty() && input.navigate_after != "10m" {
        let mut buffer = String::new();
        buffer.push_str("navigate_after(");
        let arg0 = format!("{:?}.to_string()", input.navigate_after);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
