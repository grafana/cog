namespace Grafana.Foundation.Sandbox;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder(string title)
    {
        this.@internal = new SomeStruct();
        this.@internal.Title = title;
    }

    public SomeStructBuilder Title(string title)
    {
        this.@internal.Title = title;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
