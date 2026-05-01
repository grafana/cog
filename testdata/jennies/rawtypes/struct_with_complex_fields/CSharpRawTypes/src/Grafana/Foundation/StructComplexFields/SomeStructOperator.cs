namespace Grafana.Foundation.StructComplexFields;

using System.Runtime.Serialization;

public enum SomeStructOperator
{
    [EnumMember(Value = ">")]
    GreaterThan,
    [EnumMember(Value = "<")]
    LessThan
}
