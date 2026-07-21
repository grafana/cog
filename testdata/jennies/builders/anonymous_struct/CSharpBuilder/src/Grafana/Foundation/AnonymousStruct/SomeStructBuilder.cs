namespace Grafana.Foundation.AnonymousStruct;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public SomeStructBuilder Time(object time)
    {
        this.@internal.Time = time;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
