namespace Grafana.Foundation.Sandbox;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public SomeStructBuilder Editable()
    {
        this.@internal.Editable = true;
        return this;
    }

    public SomeStructBuilder ReadonlyArg()
    {
        this.@internal.Editable = false;
        return this;
    }

    public SomeStructBuilder AutoRefresh()
    {
        this.@internal.AutoRefresh = true;
        return this;
    }

    public SomeStructBuilder NoAutoRefresh()
    {
        this.@internal.AutoRefresh = false;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
