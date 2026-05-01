namespace Grafana.Foundation.NullableMapAssignment;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public SomeStructBuilder Config(Dictionary<string, string> config)
    {
        this.@internal.Config = config;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
