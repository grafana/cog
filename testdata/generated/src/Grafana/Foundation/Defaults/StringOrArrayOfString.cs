// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Defaults;

using System.Collections.Generic;
using System.Text.Json.Serialization;

[JsonConverter(typeof(StringOrArrayOfStringJsonConverter))]
public class StringOrArrayOfString
{
    [JsonPropertyName("String")]
    public string String;
    [JsonPropertyName("ArrayOfString")]
    public List<string> ArrayOfString;

    public StringOrArrayOfString()
    {
    }

    public StringOrArrayOfString(string stringArg, List<string> arrayOfString)
    {
        this.String = stringArg;
        this.ArrayOfString = arrayOfString;
    }
}
