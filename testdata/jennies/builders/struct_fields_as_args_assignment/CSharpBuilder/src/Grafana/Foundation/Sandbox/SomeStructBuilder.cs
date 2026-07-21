namespace Grafana.Foundation.Sandbox;


public class SomeStructBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public SomeStructBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public SomeStructBuilder Time(string from,string to)
    {
        if (this.@internal.Time == null)
        {
            this.@internal.Time = new object();
        }
        this.@internal.Time.From = from;
        this.@internal.Time.To = to;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
