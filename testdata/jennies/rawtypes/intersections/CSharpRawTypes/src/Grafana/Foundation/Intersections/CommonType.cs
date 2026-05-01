namespace Grafana.Foundation.Intersections;

using System.Runtime.Serialization;

public enum CommonType
{
    [EnumMember(Value = "counter")]
    Counter,
    [EnumMember(Value = "gauge")]
    Gauge
}
