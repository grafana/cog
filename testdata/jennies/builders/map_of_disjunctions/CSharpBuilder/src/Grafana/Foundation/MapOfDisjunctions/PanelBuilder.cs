namespace Grafana.Foundation.MapOfDisjunctions;


public class PanelBuilder : Cog.IBuilder<Panel>
{
    protected readonly Panel @internal;

    public PanelBuilder()
    {
        this.@internal = new Panel();
        this.@internal.Kind = "Panel";
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
