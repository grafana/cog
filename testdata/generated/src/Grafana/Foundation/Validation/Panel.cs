// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Validation;

using System.Text.Json.Serialization;

public class Panel
{
    [JsonPropertyName("title")]
    public string Title;

    public Panel()
    {
        this.Title = "";
    }

    public Panel(string title)
    {
        this.Title = title;
    }
}
