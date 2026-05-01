namespace Grafana.Foundation.Dashboard;


public class PanelBuilder : Cog.IBuilder<Panel>
{
    protected readonly Panel @internal;

    public PanelBuilder()
    {
        this.@internal = new Panel();
    }

    public PanelBuilder OnlyFromThisDashboard(bool onlyFromThisDashboard)
    {
        this.@internal.OnlyFromThisDashboard = onlyFromThisDashboard;
        return this;
    }

    public Panel Build()
    {
        return this.@internal;
    }
}
