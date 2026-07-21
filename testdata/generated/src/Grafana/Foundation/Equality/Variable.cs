// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Equality;

using System.Text.Json.Serialization;

public class Variable
{
    [JsonPropertyName("name")]
    public string Name;

    public Variable()
    {
        this.Name = "";
    }

    public Variable(string name)
    {
        this.Name = name;
    }
}
