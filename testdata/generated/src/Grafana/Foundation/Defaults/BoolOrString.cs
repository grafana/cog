// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Defaults;

using System.Text.Json.Serialization;

[JsonConverter(typeof(BoolOrStringJsonConverter))]
public class BoolOrString
{
    [JsonPropertyName("Bool")]
    public bool? Bool;
    [JsonPropertyName("String")]
    public string String;

    public BoolOrString()
    {
    }

    public BoolOrString(bool? boolArg, string stringArg)
    {
        this.Bool = boolArg;
        this.String = stringArg;
    }
}
