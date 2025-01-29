package panelbuilder;

import java.util.List;

public class PanelBuilder<T extends PanelBuilder<T>> implements cog.Builder<Options> {
    protected final Options internal;
    
    public PanelBuilder() {
        this.internal = new Options();
    }
    public T onlyFromThisDashboard(Boolean onlyFromThisDashboard) {
        this.internal.onlyFromThisDashboard = onlyFromThisDashboard;
        return (T) this;
    }
    
    public T onlyInTimeRange(Boolean onlyInTimeRange) {
        this.internal.onlyInTimeRange = onlyInTimeRange;
        return (T) this;
    }
    
    public T tags(List<String> tags) {
        this.internal.tags = tags;
        return (T) this;
    }
    
    public T limit(Integer limit) {
        this.internal.limit = limit;
        return (T) this;
    }
    
    public T showUser(Boolean showUser) {
        this.internal.showUser = showUser;
        return (T) this;
    }
    
    public T showTime(Boolean showTime) {
        this.internal.showTime = showTime;
        return (T) this;
    }
    
    public T showTags(Boolean showTags) {
        this.internal.showTags = showTags;
        return (T) this;
    }
    
    public T navigateToPanel(Boolean navigateToPanel) {
        this.internal.navigateToPanel = navigateToPanel;
        return (T) this;
    }
    
    public T navigateBefore(String navigateBefore) {
        this.internal.navigateBefore = navigateBefore;
        return (T) this;
    }
    
    public T navigateAfter(String navigateAfter) {
        this.internal.navigateAfter = navigateAfter;
        return (T) this;
    }
    public Options build() {
        return this.internal;
    }
}
