use crate::types::builder_delegation;

/// dashboard_link_converter accepts a `DashboardLink` object and generates the Rust code to build this object using builders.
pub fn dashboard_link_converter(input: &builder_delegation::DashboardLink) -> String {
    let mut calls: Vec<String> =
        vec!["builder_delegation::DashboardLinkBuilder::new()".to_string()];
    if !input.title.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("title(");
        let arg0 = format!("{:?}.to_string()", input.title);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.url.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("url(");
        let arg0 = format!("{:?}.to_string()", input.url);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}

/// dashboard_converter accepts a `Dashboard` object and generates the Rust code to build this object using builders.
pub fn dashboard_converter(input: &builder_delegation::Dashboard) -> String {
    let mut calls: Vec<String> = vec!["builder_delegation::DashboardBuilder::new()".to_string()];
    {
        let mut buffer = String::new();
        buffer.push_str("id(");
        let arg0 = format!("{:?}", input.id);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.title.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("title(");
        let arg0 = format!("{:?}.to_string()", input.title);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.links.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("links(");
        let mut tmp_arg0: Vec<String> = Vec::new();
        for arg1 in &input.links {
            let tmp_links_arg1 = dashboard_link_converter(arg1);
            tmp_arg0.push(tmp_links_arg1);
        }
        let arg0 = format!("vec![{}]", tmp_arg0.join(", "));
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    if !input.links_of_links.is_empty() {
        let mut buffer = String::new();
        buffer.push_str("links_of_links(");
        let mut tmp_arg0: Vec<String> = Vec::new();
        for arg1 in &input.links_of_links {
            let mut tmp_tmp_links_of_links_arg1: Vec<String> = Vec::new();
            for arg1_value in arg1 {
                let tmp_arg1_arg1_value = dashboard_link_converter(arg1_value);
                tmp_tmp_links_of_links_arg1.push(tmp_arg1_arg1_value);
            }
            let tmp_links_of_links_arg1 =
                format!("vec![{}]", tmp_tmp_links_of_links_arg1.join(", "));
            tmp_arg0.push(tmp_links_of_links_arg1);
        }
        let arg0 = format!("vec![{}]", tmp_arg0.join(", "));
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }
    {
        let mut buffer = String::new();
        buffer.push_str("single_link(");
        let arg0 = dashboard_link_converter(&input.single_link);
        buffer.push_str(&arg0);
        buffer.push(')');
        calls.push(buffer);
    }

    calls.join("\n    .")
}
