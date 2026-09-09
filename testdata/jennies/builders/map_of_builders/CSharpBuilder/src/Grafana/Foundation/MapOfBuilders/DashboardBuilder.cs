namespace Grafana.Foundation.MapOfBuilders;


public class DashboardBuilder : Cog.IBuilder<Dashboard>
{
    protected readonly Dashboard @internal;

    public DashboardBuilder()
    {
        this.@internal = new Dashboard();
    }

    public DashboardBuilder Panels(Dictionary<string, Cog.IBuilder<Panel>> panels)
    {
        Dictionary<string, Panel> panelsResources = new Dictionary<string, Panel>();
        foreach (var entry1 in panels)
        {
                Panel panelsDepth1 = entry1.Value.Build();
                panelsResources[entry1.Key] = panelsDepth1;
        }
        this.@internal.Panels = panelsResources;
        return this;
    }

    public Dashboard Build()
    {
        return this.@internal;
    }
}
