namespace Grafana.Foundation.Properties;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;
    private string someBuilderProperty;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
        this.someBuilderProperty = "";
    }

    public SomeStructBuilder Id(long id)
    {
        this.@internal.Id = id;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
