namespace Grafana.Foundation.Enums;

using System.Runtime.Serialization;

// This is a very interesting string enum.
public enum Operator
{
    [EnumMember(Value = ">")]
    GreaterThan,
    [EnumMember(Value = "<")]
    LessThan
}
