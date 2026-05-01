namespace Grafana.Foundation.Panelbuilder;


public class PanelBuilder : Cog.IBuilder<Options>
{
    protected readonly Options @internal;

    public PanelBuilder()
    {
        this.@internal = new Options();
    }

    public PanelBuilder OnlyFromThisDashboard(bool onlyFromThisDashboard)
    {
        this.@internal.OnlyFromThisDashboard = onlyFromThisDashboard;
        return this;
    }

    public PanelBuilder OnlyInTimeRange(bool onlyInTimeRange)
    {
        this.@internal.OnlyInTimeRange = onlyInTimeRange;
        return this;
    }

    public PanelBuilder Tags(List<string> tags)
    {
        this.@internal.Tags = tags;
        return this;
    }

    public PanelBuilder Limit(uint limit)
    {
        this.@internal.Limit = limit;
        return this;
    }

    public PanelBuilder ShowUser(bool showUser)
    {
        this.@internal.ShowUser = showUser;
        return this;
    }

    public PanelBuilder ShowTime(bool showTime)
    {
        this.@internal.ShowTime = showTime;
        return this;
    }

    public PanelBuilder ShowTags(bool showTags)
    {
        this.@internal.ShowTags = showTags;
        return this;
    }

    public PanelBuilder NavigateToPanel(bool navigateToPanel)
    {
        this.@internal.NavigateToPanel = navigateToPanel;
        return this;
    }

    public PanelBuilder NavigateBefore(string navigateBefore)
    {
        this.@internal.NavigateBefore = navigateBefore;
        return this;
    }

    public PanelBuilder NavigateAfter(string navigateAfter)
    {
        this.@internal.NavigateAfter = navigateAfter;
        return this;
    }

    public Options Build()
    {
        return this.@internal;
    }
}
