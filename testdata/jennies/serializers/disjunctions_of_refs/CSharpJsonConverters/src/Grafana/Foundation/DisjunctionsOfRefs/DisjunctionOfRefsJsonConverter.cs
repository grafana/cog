namespace Grafana.Foundation.DisjunctionsOfRefs;

using System.Text.Json;
using System.Text.Json.Serialization;

public sealed class DisjunctionOfRefsJsonConverter : JsonConverter<DisjunctionOfRefs>
{
    public override DisjunctionOfRefs Read(ref Utf8JsonReader reader, System.Type typeToConvert, JsonSerializerOptions options)
    {
        using var doc = JsonDocument.ParseValue(ref reader);
        var root = doc.RootElement;
        if (!root.TryGetProperty("type", out var discProp))
        {
            throw new JsonException("missing discriminator property 'type' for DisjunctionOfRefs");
        }

        var result = new DisjunctionOfRefs();
        var raw = root.GetRawText();
        switch (discProp.GetString())
        {
            case "A":
                result.MyRefA = JsonSerializer.Deserialize<MyRefA>(raw, options);
                break;
            case "B":
                result.MyRefB = JsonSerializer.Deserialize<MyRefB>(raw, options);
                break;
            default:
                throw new JsonException("unknown discriminator value '" + discProp.GetString() + "' for DisjunctionOfRefs");
        }
        return result;
    }

    public override void Write(Utf8JsonWriter writer, DisjunctionOfRefs value, JsonSerializerOptions options)
    {
        if (value.MyRefA != null) { JsonSerializer.Serialize(writer, value.MyRefA, options); }
        else if (value.MyRefB != null) { JsonSerializer.Serialize(writer, value.MyRefB, options); }
        else
        {
            writer.WriteNullValue();
        }
    }
}
