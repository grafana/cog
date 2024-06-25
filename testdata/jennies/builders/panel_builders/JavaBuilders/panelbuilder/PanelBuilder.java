package panelbuilder;

import java.util.List;
import dashboard.Panel;

public class PanelBuilder implements cog.Builder<Panel> {
    private Panel internal;

    public PanelBuilder() {
        this.internal = new Panel();
        this.OnlyFromThisDashboard(false);
        this.OnlyInTimeRange(false);
        this.Limit(10);
        this.ShowUser(true);
        this.ShowTime(true);
        this.ShowTags(true);
        this.NavigateToPanel(true);
        this.NavigateBefore("10m");
        this.NavigateAfter("10m");
    }
    public PanelBuilder OnlyFromThisDashboard(Boolean onlyFromThisDashboard) {
    this.internal.onlyFromThisDashboard = onlyFromThisDashboard;
        return this;
    }
    public PanelBuilder OnlyInTimeRange(Boolean onlyInTimeRange) {
    this.internal.onlyInTimeRange = onlyInTimeRange;
        return this;
    }
    public PanelBuilder Tags(List<String> tags) {
    this.internal.tags = tags;
        return this;
    }
    public PanelBuilder Limit(Integer limit) {
    this.internal.limit = limit;
        return this;
    }
    public PanelBuilder ShowUser(Boolean showUser) {
    this.internal.showUser = showUser;
        return this;
    }
    public PanelBuilder ShowTime(Boolean showTime) {
    this.internal.showTime = showTime;
        return this;
    }
    public PanelBuilder ShowTags(Boolean showTags) {
    this.internal.showTags = showTags;
        return this;
    }
    public PanelBuilder NavigateToPanel(Boolean navigateToPanel) {
    this.internal.navigateToPanel = navigateToPanel;
        return this;
    }
    public PanelBuilder NavigateBefore(String navigateBefore) {
    this.internal.navigateBefore = navigateBefore;
        return this;
    }
    public PanelBuilder NavigateAfter(String navigateAfter) {
    this.internal.navigateAfter = navigateAfter;
        return this;
    }
    
    public Panel build() {
        return this.internal;
    }
}
