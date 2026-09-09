// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Equality;

using System.Collections.Generic;
using System.Text.Json.Serialization;

public class Maps
{
    [JsonPropertyName("ints")]
    public Dictionary<string, long> Ints;
    [JsonPropertyName("strings")]
    public Dictionary<string, string> Strings;
    [JsonPropertyName("refs")]
    public Dictionary<string, Variable> Refs;
    [JsonPropertyName("anonymousStructs")]
    public Dictionary<string, EqualityMapsAnonymousStructs> AnonymousStructs;
    [JsonPropertyName("stringToAny")]
    public Dictionary<string, object> StringToAny;

    public Maps()
    {
        this.Ints = new Dictionary<string, long>();
        this.Strings = new Dictionary<string, string>();
        this.Refs = new Dictionary<string, Variable>();
        this.AnonymousStructs = new Dictionary<string, EqualityMapsAnonymousStructs>();
        this.StringToAny = new Dictionary<string, object>();
    }

    public Maps(Dictionary<string, long> ints, Dictionary<string, string> strings, Dictionary<string, Variable> refs, Dictionary<string, EqualityMapsAnonymousStructs> anonymousStructs, Dictionary<string, object> stringToAny)
    {
        this.Ints = ints;
        this.Strings = strings;
        this.Refs = refs;
        this.AnonymousStructs = anonymousStructs;
        this.StringToAny = stringToAny;
    }
}
