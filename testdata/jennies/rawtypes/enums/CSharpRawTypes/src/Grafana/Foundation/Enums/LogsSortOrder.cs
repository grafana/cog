namespace Grafana.Foundation.Enums;

using System.Runtime.Serialization;

public enum LogsSortOrder
{
    [EnumMember(Value = "time_asc")]
    Asc,
    [EnumMember(Value = "time_desc")]
    Desc
}
