package panelbuilder;

import java.util.List;
import dashboard.Panel;

public class PanelBuilder {
    private Panel internal;

    public PanelBuilder() {
        this.internal = new Panel();
    }
    public PanelBuilder setOnlyFromThisDashboard(Boolean onlyFromThisDashboard) {
    this.internal.onlyFromThisDashboard = onlyFromThisDashboard;
        return this;
    }
    public PanelBuilder setOnlyInTimeRange(Boolean onlyInTimeRange) {
    this.internal.onlyInTimeRange = onlyInTimeRange;
        return this;
    }
    public PanelBuilder setTags(List<String> tags) {
    this.internal.tags = tags;
        return this;
    }
    public PanelBuilder setLimit(Integer limit) {
    this.internal.limit = limit;
        return this;
    }
    public PanelBuilder setShowUser(Boolean showUser) {
    this.internal.showUser = showUser;
        return this;
    }
    public PanelBuilder setShowTime(Boolean showTime) {
    this.internal.showTime = showTime;
        return this;
    }
    public PanelBuilder setShowTags(Boolean showTags) {
    this.internal.showTags = showTags;
        return this;
    }
    public PanelBuilder setNavigateToPanel(Boolean navigateToPanel) {
    this.internal.navigateToPanel = navigateToPanel;
        return this;
    }
    public PanelBuilder setNavigateBefore(String navigateBefore) {
    this.internal.navigateBefore = navigateBefore;
        return this;
    }
    public PanelBuilder setNavigateAfter(String navigateAfter) {
    this.internal.navigateAfter = navigateAfter;
        return this;
    }
    
    public Panel build() {
        return this.internal;
    }
}
