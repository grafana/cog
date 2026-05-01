namespace Grafana.Foundation.BuilderDelegation;


public class DashboardBuilder : Cog.IBuilder<Dashboard>
{
    protected readonly Dashboard @internal;

    public DashboardBuilder()
    {
        this.@internal = new Dashboard();
    }

    public DashboardBuilder Id(long id)
    {
        this.@internal.Id = id;
        return this;
    }

    public DashboardBuilder Title(string title)
    {
        this.@internal.Title = title;
        return this;
    }

    public DashboardBuilder Links(List<Cog.IBuilder<DashboardLink>> links)
    {
        List<DashboardLink> linksResources = new List<DashboardLink>();
        foreach (Cog.IBuilder<DashboardLink> r1 in links)
        {
                DashboardLink linksDepth1 = r1.Build();
                linksResources.Add(linksDepth1);
        }
        this.@internal.Links = linksResources;
        return this;
    }

    public DashboardBuilder LinksOfLinks(List<List<Cog.IBuilder<DashboardLink>>> linksOfLinks)
    {
        List<List<DashboardLink>> linksOfLinksResources = new List<List<DashboardLink>>();
        foreach (List<Cog.IBuilder<DashboardLink>> r1 in linksOfLinks)
        {
                List<DashboardLink> linksOfLinksDepth1 = new List<DashboardLink>();
        foreach (Cog.IBuilder<DashboardLink> r2 in r1)
        {
                DashboardLink linksOfLinksDepth2 = r2.Build();
                linksOfLinksDepth1.Add(linksOfLinksDepth2);
        }

                linksOfLinksResources.Add(linksOfLinksDepth1);
        }
        this.@internal.LinksOfLinks = linksOfLinksResources;
        return this;
    }

    public DashboardBuilder SingleLink(Cog.IBuilder<DashboardLink> singleLink)
    {
        DashboardLink singleLinkResource = singleLink.Build();
        this.@internal.SingleLink = singleLinkResource;
        return this;
    }

    public Dashboard Build()
    {
        return this.@internal;
    }
}
