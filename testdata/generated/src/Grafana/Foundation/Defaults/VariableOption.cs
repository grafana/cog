// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Defaults;

using System.Text.Json.Serialization;

public class VariableOption
{
    [JsonPropertyName("selected")]
    public BoolOrString Selected;
    [JsonPropertyName("text")]
    public StringOrArrayOfString Text;
    [JsonPropertyName("value")]
    public StringOrArrayOfString Value;

    public VariableOption()
    {
        this.Text = new StringOrArrayOfString();
        this.Value = new StringOrArrayOfString();
    }

    public VariableOption(BoolOrString selected, StringOrArrayOfString text, StringOrArrayOfString value)
    {
        this.Selected = selected;
        this.Text = text;
        this.Value = value;
    }
}
