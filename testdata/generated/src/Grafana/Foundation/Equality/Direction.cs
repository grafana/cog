// Code generated - EDITING IS FUTILE. DO NOT EDIT.

namespace Grafana.Foundation.Equality;

using System.Text.Json.Serialization;

[JsonConverter(typeof(JsonStringEnumConverter<Direction>))]
public enum Direction
{
    [JsonStringEnumMemberName("top")]
    Top,
    [JsonStringEnumMemberName("bottom")]
    Bottom,
    [JsonStringEnumMemberName("left")]
    Left,
    [JsonStringEnumMemberName("right")]
    Right
}
