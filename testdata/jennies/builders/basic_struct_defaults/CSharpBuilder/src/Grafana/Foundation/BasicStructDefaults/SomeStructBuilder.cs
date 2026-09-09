namespace Grafana.Foundation.BasicStructDefaults;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public SomeStructBuilder Id(long id)
    {
        this.@internal.Id = id;
        return this;
    }

    public SomeStructBuilder Uid(string uid)
    {
        this.@internal.Uid = uid;
        return this;
    }

    public SomeStructBuilder Tags(List<string> tags)
    {
        this.@internal.Tags = tags;
        return this;
    }

    public SomeStructBuilder LiveNow(bool liveNow)
    {
        this.@internal.LiveNow = liveNow;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
