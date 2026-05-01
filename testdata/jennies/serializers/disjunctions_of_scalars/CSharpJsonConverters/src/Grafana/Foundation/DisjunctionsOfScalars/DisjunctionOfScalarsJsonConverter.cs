namespace Grafana.Foundation.DisjunctionsOfScalars;

using System.Text.Json;
using System.Text.Json.Serialization;

public sealed class DisjunctionOfScalarsJsonConverter : JsonConverter<DisjunctionOfScalars>
{
    public override DisjunctionOfScalars Read(ref Utf8JsonReader reader, System.Type typeToConvert, JsonSerializerOptions options)
    {
        var result = new DisjunctionOfScalars();
        switch (reader.TokenType)
        {
            case JsonTokenType.String:
                result.String = reader.GetString();
                break;
            case JsonTokenType.True:
            case JsonTokenType.False:
                result.Bool = reader.GetBoolean();
                break;
            case JsonTokenType.StartArray:
                result.ArrayOfString = JsonSerializer.Deserialize<List<string>>(ref reader, options);
                break;
            case JsonTokenType.Number:
                if (reader.TryGetInt64(out _)) { result.Int64 = reader.GetInt64(); }
                else { result.Float32 = reader.GetSingle(); }
                break;
            default:
                throw new JsonException("unexpected token while deserialising DisjunctionOfScalars: " + reader.TokenType);
        }
        return result;
    }

    public override void Write(Utf8JsonWriter writer, DisjunctionOfScalars value, JsonSerializerOptions options)
    {
        if (value.String != null) { JsonSerializer.Serialize(writer, value.String, options); }
        else if (value.Bool != null) { JsonSerializer.Serialize(writer, value.Bool.Value, options); }
        else if (value.ArrayOfString != null) { JsonSerializer.Serialize(writer, value.ArrayOfString, options); }
        else if (value.Int64 != null) { JsonSerializer.Serialize(writer, value.Int64.Value, options); }
        else if (value.Float32 != null) { JsonSerializer.Serialize(writer, value.Float32.Value, options); }
        else
        {
            writer.WriteNullValue();
        }
    }
}
