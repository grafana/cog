namespace Grafana.Foundation.Intersections;

using System.Runtime.Serialization;

public enum CommonContains
{
    [EnumMember(Value = "default")]
    Default,
    [EnumMember(Value = "time")]
    Time
}
