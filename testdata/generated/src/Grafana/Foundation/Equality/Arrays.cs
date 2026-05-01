// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Equality;

using System.Collections.Generic;
using System.Text.Json.Serialization;

public class Arrays
{
    [JsonPropertyName("ints")]
    public List<long> Ints;
    [JsonPropertyName("strings")]
    public List<string> Strings;
    [JsonPropertyName("arrayOfArray")]
    public List<List<string>> ArrayOfArray;
    [JsonPropertyName("refs")]
    public List<Variable> Refs;
    [JsonPropertyName("anonymousStructs")]
    public List<EqualityArraysAnonymousStructs> AnonymousStructs;
    [JsonPropertyName("arrayOfAny")]
    public List<object> ArrayOfAny;

    public Arrays()
    {
        this.Ints = new List<long>();
        this.Strings = new List<string>();
        this.ArrayOfArray = new List<List<string>>();
        this.Refs = new List<Variable>();
        this.AnonymousStructs = new List<EqualityArraysAnonymousStructs>();
        this.ArrayOfAny = new List<object>();
    }

    public Arrays(List<long> ints, List<string> strings, List<List<string>> arrayOfArray, List<Variable> refs, List<EqualityArraysAnonymousStructs> anonymousStructs, List<object> arrayOfAny)
    {
        this.Ints = ints;
        this.Strings = strings;
        this.ArrayOfArray = arrayOfArray;
        this.Refs = refs;
        this.AnonymousStructs = anonymousStructs;
        this.ArrayOfAny = arrayOfAny;
    }
}
