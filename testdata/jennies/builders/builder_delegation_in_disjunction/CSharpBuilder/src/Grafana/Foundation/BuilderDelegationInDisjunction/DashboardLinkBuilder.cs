namespace Grafana.Foundation.BuilderDelegationInDisjunction;


public class DashboardLinkBuilder : Cog.IBuilder<DashboardLink>
{
    protected readonly DashboardLink @internal;

    public DashboardLinkBuilder()
    {
        this.@internal = new DashboardLink();
    }

    public DashboardLinkBuilder Title(string title)
    {
        this.@internal.Title = title;
        return this;
    }

    public DashboardLinkBuilder Url(string url)
    {
        this.@internal.Url = url;
        return this;
    }

    public DashboardLink Build()
    {
        return this.@internal;
    }
}
