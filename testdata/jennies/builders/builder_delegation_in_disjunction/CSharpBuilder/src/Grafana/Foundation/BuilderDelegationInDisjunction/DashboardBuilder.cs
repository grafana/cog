namespace Grafana.Foundation.BuilderDelegationInDisjunction;


public class DashboardBuilder : Cog.IBuilder<Dashboard>
{
    protected readonly Dashboard @internal;

    public DashboardBuilder()
    {
        this.@internal = new Dashboard();
    }

    public DashboardBuilder SingleLinkOrString(Cog.IBuilder<object> singleLinkOrString)
    {
        object singleLinkOrStringResource = singleLinkOrString.Build();
        this.@internal.SingleLinkOrString = singleLinkOrStringResource;
        return this;
    }

    public DashboardBuilder LinksOrStrings(List<Cog.IBuilder<object>> linksOrStrings)
    {
        List<object> linksOrStringsResources = new List<object>();
        foreach (Cog.IBuilder<object> r1 in linksOrStrings)
        {
                object linksOrStringsDepth1 = r1.Build();
                linksOrStringsResources.Add(linksOrStringsDepth1);
        }
        this.@internal.LinksOrStrings = linksOrStringsResources;
        return this;
    }

    public DashboardBuilder DisjunctionOfBuilders(Cog.IBuilder<object> disjunctionOfBuilders)
    {
        object disjunctionOfBuildersResource = disjunctionOfBuilders.Build();
        this.@internal.DisjunctionOfBuilders = disjunctionOfBuildersResource;
        return this;
    }

    public Dashboard Build()
    {
        return this.@internal;
    }
}
