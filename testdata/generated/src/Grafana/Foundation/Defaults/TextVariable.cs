// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Defaults;

using System.Text.Json.Serialization;

public class TextVariable
{
    [JsonPropertyName("name")]
    public string Name;
    [JsonPropertyName("current")]
    public VariableOption Current;
    [JsonPropertyName("skipUrlSync")]
    public bool SkipUrlSync;

    public TextVariable()
    {
        this.Name = "";
        this.Current = new VariableOption();
        this.SkipUrlSync = false;
    }

    public TextVariable(string name, VariableOption current, bool skipUrlSync)
    {
        this.Name = name;
        this.Current = current;
        this.SkipUrlSync = skipUrlSync;
    }
}
