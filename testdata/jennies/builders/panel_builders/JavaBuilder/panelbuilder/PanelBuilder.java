package panelbuilder;

import java.util.List;

public class PanelBuilder implements cog.Builder<Options> {
    protected final Options internal;
    
    public PanelBuilder() {
        this.internal = new Options();
    }
    public PanelBuilder onlyFromThisDashboard(Boolean onlyFromThisDashboard) {
        this.internal.onlyFromThisDashboard = onlyFromThisDashboard;
        return this;
    }
    
    public PanelBuilder onlyInTimeRange(Boolean onlyInTimeRange) {
        this.internal.onlyInTimeRange = onlyInTimeRange;
        return this;
    }
    
    public PanelBuilder tags(List<String> tags) {
        this.internal.tags = tags;
        return this;
    }
    
    public PanelBuilder limit(Integer limit) {
        this.internal.limit = limit;
        return this;
    }
    
    public PanelBuilder showUser(Boolean showUser) {
        this.internal.showUser = showUser;
        return this;
    }
    
    public PanelBuilder showTime(Boolean showTime) {
        this.internal.showTime = showTime;
        return this;
    }
    
    public PanelBuilder showTags(Boolean showTags) {
        this.internal.showTags = showTags;
        return this;
    }
    
    public PanelBuilder navigateToPanel(Boolean navigateToPanel) {
        this.internal.navigateToPanel = navigateToPanel;
        return this;
    }
    
    public PanelBuilder navigateBefore(String navigateBefore) {
        this.internal.navigateBefore = navigateBefore;
        return this;
    }
    
    public PanelBuilder navigateAfter(String navigateAfter) {
        this.internal.navigateAfter = navigateAfter;
        return this;
    }
    public Options build() {
        return this.internal;
    }
}
