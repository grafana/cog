// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Equality;

using System.Text.Json.Serialization;

public class Container
{
    [JsonPropertyName("stringField")]
    public string StringField;
    [JsonPropertyName("intField")]
    public long IntField;
    [JsonPropertyName("enumField")]
    public Direction EnumField;
    [JsonPropertyName("refField")]
    public Variable RefField;

    public Container()
    {
        this.StringField = "";
        this.IntField = 0L;
        this.EnumField = Direction.Top;
        this.RefField = new Variable();
    }

    public Container(string stringField, long intField, Direction enumField, Variable refField)
    {
        this.StringField = stringField;
        this.IntField = intField;
        this.EnumField = enumField;
        this.RefField = refField;
    }
}
