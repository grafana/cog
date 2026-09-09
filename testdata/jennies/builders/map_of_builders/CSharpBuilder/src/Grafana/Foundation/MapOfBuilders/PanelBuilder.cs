namespace Grafana.Foundation.MapOfBuilders;


public class PanelBuilder : Cog.IBuilder<Panel>
{
    protected readonly Panel @internal;

    public PanelBuilder()
    {
        this.@internal = new Panel();
    }

    public PanelBuilder Title(string title)
    {
        this.@internal.Title = title;
        return this;
    }

    public Panel Build()
    {
        return this.@internal;
    }
}
