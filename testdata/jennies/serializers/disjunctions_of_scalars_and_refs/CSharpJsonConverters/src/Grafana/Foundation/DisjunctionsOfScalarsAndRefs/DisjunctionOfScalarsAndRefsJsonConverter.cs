namespace Grafana.Foundation.DisjunctionsOfScalarsAndRefs;

using System.Text.Json;
using System.Text.Json.Serialization;

public sealed class DisjunctionOfScalarsAndRefsJsonConverter : JsonConverter<DisjunctionOfScalarsAndRefs>
{
    public override DisjunctionOfScalarsAndRefs Read(ref Utf8JsonReader reader, System.Type typeToConvert, JsonSerializerOptions options)
    {
        using var doc = JsonDocument.ParseValue(ref reader);
        var root = doc.RootElement;
        var raw = root.GetRawText();
        var result = new DisjunctionOfScalarsAndRefs();

        switch (root.ValueKind)
        {
            case JsonValueKind.String:
                result.String = root.GetString();
                break;
            case JsonValueKind.True:
            case JsonValueKind.False:
                result.Bool = root.GetBoolean();
                break;
            case JsonValueKind.Array:
                result.ArrayOfString = JsonSerializer.Deserialize<List<string>>(raw, options);
                break;
            case JsonValueKind.Object:
                if (TryDeserialize<MyRefA>(raw, options, out var __MyRefA)) { result.MyRefA = __MyRefA; break; }
                if (TryDeserialize<MyRefB>(raw, options, out var __MyRefB)) { result.MyRefB = __MyRefB; break; }
                result.Any = JsonSerializer.Deserialize<JsonElement>(raw, options);
                break;
            default:
                throw new JsonException("unexpected token while deserialising DisjunctionOfScalarsAndRefs: " + root.ValueKind);
        }

        return result;
    }

    public override void Write(Utf8JsonWriter writer, DisjunctionOfScalarsAndRefs value, JsonSerializerOptions options)
    {
        if (value.String != null) { JsonSerializer.Serialize(writer, value.String, options); }
        else if (value.Bool != null) { JsonSerializer.Serialize(writer, value.Bool.Value, options); }
        else if (value.ArrayOfString != null) { JsonSerializer.Serialize(writer, value.ArrayOfString, options); }
        else if (value.MyRefA != null) { JsonSerializer.Serialize(writer, value.MyRefA, options); }
        else if (value.MyRefB != null) { JsonSerializer.Serialize(writer, value.MyRefB, options); }
        else if (value.Any != null) { JsonSerializer.Serialize(writer, value.Any, options); }
        else
        {
            writer.WriteNullValue();
        }
    }

    private static bool TryDeserialize<T>(string raw, JsonSerializerOptions options, out T? value)
    {
        try
        {
            value = JsonSerializer.Deserialize<T>(raw, options);
            return value != null;
        }
        catch (JsonException)
        {
            value = default;
            return false;
        }
    }
}
