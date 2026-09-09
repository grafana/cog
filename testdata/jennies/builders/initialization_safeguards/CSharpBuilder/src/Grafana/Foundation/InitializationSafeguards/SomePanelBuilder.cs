namespace Grafana.Foundation.InitializationSafeguards;


public class SomePanelBuilder : Cog.IBuilder<SomePanel>
{
    protected readonly SomePanel @internal;

    public SomePanelBuilder()
    {
        this.@internal = new SomePanel();
    }

    public SomePanelBuilder Title(string title)
    {
        this.@internal.Title = title;
        return this;
    }

    public SomePanelBuilder ShowLegend(bool show)
    {
        if (this.@internal.Options == null)
        {
            this.@internal.Options = new Options();
        }
        if (this.@internal.Options.Legend == null)
        {
            this.@internal.Options.Legend = new LegendOptions();
        }
        this.@internal.Options.Legend.Show = show;
        return this;
    }

    public SomePanel Build()
    {
        return this.@internal;
    }
}
