namespace Grafana.Foundation.KnownAny;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public SomeStructBuilder Title(string title)
    {
        if (this.@internal.Config == null)
        {
            this.@internal.Config = new Config();
        }
        ((Config) this.@internal.Config).Title = title;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
