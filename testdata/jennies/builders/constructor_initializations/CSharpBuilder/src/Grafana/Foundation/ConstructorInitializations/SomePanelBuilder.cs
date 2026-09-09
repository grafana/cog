namespace Grafana.Foundation.ConstructorInitializations;


public class SomePanelBuilder : Cog.IBuilder<SomePanel>
{
    protected readonly SomePanel @internal;

    public SomePanelBuilder()
    {
        this.@internal = new SomePanel();
        this.@internal.Type = "panel_type";
        this.@internal.Cursor = CursorMode.Tooltip;
    }

    public SomePanelBuilder Title(string title)
    {
        this.@internal.Title = title;
        return this;
    }

    public SomePanel Build()
    {
        return this.@internal;
    }
}
