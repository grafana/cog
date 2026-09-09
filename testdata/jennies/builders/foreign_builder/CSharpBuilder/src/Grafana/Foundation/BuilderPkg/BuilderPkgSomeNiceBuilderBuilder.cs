namespace Grafana.Foundation.BuilderPkg;

using Grafana.Foundation.SomePkg;

public class BuilderPkgSomeNiceBuilderBuilder : Cog.IBuilder<SomeStruct>
{
    protected readonly SomeStruct @internal;

    public BuilderPkgSomeNiceBuilderBuilder()
    {
        this.@internal = new SomeStruct();
    }

    public BuilderPkgSomeNiceBuilderBuilder Title(string title)
    {
        this.@internal.Title = title;
        return this;
    }

    public SomeStruct Build()
    {
        return this.@internal;
    }
}
