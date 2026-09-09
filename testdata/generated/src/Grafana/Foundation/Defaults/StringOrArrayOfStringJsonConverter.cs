// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Defaults;

using System.Text.Json;
using System.Text.Json.Serialization;

public sealed class StringOrArrayOfStringJsonConverter : JsonConverter<StringOrArrayOfString>
{
    public override StringOrArrayOfString Read(ref Utf8JsonReader reader, System.Type typeToConvert, JsonSerializerOptions options)
    {
        var result = new StringOrArrayOfString();
        switch (reader.TokenType)
        {
            case JsonTokenType.String:
                result.String = reader.GetString();
                break;
            case JsonTokenType.StartArray:
                result.ArrayOfString = JsonSerializer.Deserialize<List<string>>(ref reader, options);
                break;
            default:
                throw new JsonException("unexpected token while deserialising StringOrArrayOfString: " + reader.TokenType);
        }
        return result;
    }

    public override void Write(Utf8JsonWriter writer, StringOrArrayOfString value, JsonSerializerOptions options)
    {
        if (value.String != null) { JsonSerializer.Serialize(writer, value.String, options); }
        else if (value.ArrayOfString != null) { JsonSerializer.Serialize(writer, value.ArrayOfString, options); }
        else
        {
            writer.WriteNullValue();
        }
    }
}
