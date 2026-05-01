namespace Grafana.Foundation.Sandbox;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public SomeStructBuilder Annotations(string key,string value)
    {
        if (this.@internal.Annotations == null)
        {
            this.@internal.Annotations = new Dictionary<string, string>();
        }
        this.@internal.Annotations[key] = value;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
