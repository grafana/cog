// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Defaults;

using System.Text.Json;
using System.Text.Json.Serialization;

public sealed class BoolOrStringJsonConverter : JsonConverter<BoolOrString>
{
    public override BoolOrString Read(ref Utf8JsonReader reader, System.Type typeToConvert, JsonSerializerOptions options)
    {
        var result = new BoolOrString();
        switch (reader.TokenType)
        {
            case JsonTokenType.True:
            case JsonTokenType.False:
                result.Bool = reader.GetBoolean();
                break;
            case JsonTokenType.String:
                result.String = reader.GetString();
                break;
            default:
                throw new JsonException("unexpected token while deserialising BoolOrString: " + reader.TokenType);
        }
        return result;
    }

    public override void Write(Utf8JsonWriter writer, BoolOrString value, JsonSerializerOptions options)
    {
        if (value.Bool != null) { JsonSerializer.Serialize(writer, value.Bool.Value, options); }
        else if (value.String != null) { JsonSerializer.Serialize(writer, value.String, options); }
        else
        {
            writer.WriteNullValue();
        }
    }
}
