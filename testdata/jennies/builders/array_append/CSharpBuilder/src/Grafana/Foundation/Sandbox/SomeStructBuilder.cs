namespace Grafana.Foundation.Sandbox;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public SomeStructBuilder Tags(string tags)
    {
        if (this.@internal.Tags == null)
        {
            this.@internal.Tags = new List<string>();
        }
        this.@internal.Tags.Add(tags);
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
