namespace Grafana.Foundation.StructWithDefaults;


public class NestedStructBuilder : Cog.IBuilder<NestedStruct>
{
    protected readonly NestedStruct @internal;

    public NestedStructBuilder()
    {
        this.@internal = new NestedStruct();
    }

    public NestedStructBuilder StringVal(string stringVal)
    {
        this.@internal.StringVal = stringVal;
        return this;
    }

    public NestedStructBuilder IntVal(long intVal)
    {
        this.@internal.IntVal = intVal;
        return this;
    }

    public NestedStruct Build()
    {
        return this.@internal;
    }
}
