// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Equality;

using System.Text.Json.Serialization;

public class Optionals
{
    [JsonPropertyName("stringField")]
    public string StringField;
    [JsonPropertyName("enumField")]
    public Direction EnumField;
    [JsonPropertyName("refField")]
    public Variable RefField;
    [JsonPropertyName("byteField")]
    public byte ByteField;

    public Optionals()
    {
    }

    public Optionals(string stringField, Direction enumField, Variable refField, byte byteField)
    {
        this.StringField = stringField;
        this.EnumField = enumField;
        this.RefField = refField;
        this.ByteField = byteField;
    }
}
