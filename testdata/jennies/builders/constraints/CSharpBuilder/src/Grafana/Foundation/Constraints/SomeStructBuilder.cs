namespace Grafana.Foundation.Constraints;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public SomeStructBuilder Id(ulong id)
    {
        if (!(id >= 5))
        {
            throw new System.ArgumentException("id must be >= 5");
        }
        if (!(id < 10))
        {
            throw new System.ArgumentException("id must be < 10");
        }
        this.@internal.Id = id;
        return this;
    }

    public SomeStructBuilder Title(string title)
    {
        if (!(title.Length >= 1))
        {
            throw new System.ArgumentException("title.Length must be >= 1");
        }
        this.@internal.Title = title;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
