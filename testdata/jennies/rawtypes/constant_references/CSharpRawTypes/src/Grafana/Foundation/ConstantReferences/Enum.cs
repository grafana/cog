namespace Grafana.Foundation.ConstantReferences;

using System.Runtime.Serialization;

public enum Enum
{
    [EnumMember(Value = "ValueA")]
    ValueA,
    [EnumMember(Value = "ValueB")]
    ValueB,
    [EnumMember(Value = "ValueC")]
    ValueC
}
