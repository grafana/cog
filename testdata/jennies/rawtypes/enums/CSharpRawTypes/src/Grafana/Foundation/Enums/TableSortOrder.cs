namespace Grafana.Foundation.Enums;

using System.Runtime.Serialization;

public enum TableSortOrder
{
    [EnumMember(Value = "asc")]
    Asc,
    [EnumMember(Value = "desc")]
    Desc
}
