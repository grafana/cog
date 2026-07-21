// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Validation;

using System.Collections.Generic;
using System.Text.Json.Serialization;

public class Dashboard
{
    [JsonPropertyName("uid")]
    public string Uid;
    [JsonPropertyName("id")]
    public long Id;
    [JsonPropertyName("title")]
    public string Title;
    [JsonPropertyName("tags")]
    public List<string> Tags;
    [JsonPropertyName("labels")]
    public Dictionary<string, string> Labels;
    [JsonPropertyName("panels")]
    public List<Panel> Panels;

    public Dashboard()
    {
        this.Title = "";
        this.Tags = new List<string>();
        this.Labels = new Dictionary<string, string>();
        this.Panels = new List<Panel>();
    }

    public Dashboard(string uid, long id, string title, List<string> tags, Dictionary<string, string> labels, List<Panel> panels)
    {
        this.Uid = uid;
        this.Id = id;
        this.Title = title;
        this.Tags = tags;
        this.Labels = labels;
        this.Panels = panels;
    }
}
