namespace Grafana.Foundation.MapOfDisjunctions;


public class DashboardBuilder : Cog.IBuilder<Dashboard>
{
    protected readonly Dashboard @internal;

    public DashboardBuilder()
    {
        this.@internal = new Dashboard();
    }

    public DashboardBuilder Panels(Dictionary<string, Cog.IBuilder<Element>> panels)
    {
        Dictionary<string, Element> panelsResources = new Dictionary<string, Element>();
        foreach (var entry1 in panels)
        {
                Element panelsDepth1 = entry1.Value.Build();
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
