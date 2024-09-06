package panelbuilder;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonSetter;
import com.fasterxml.jackson.annotation.Nulls;
import java.util.List;
import dashboard.Panel;

public class PanelBuilder implements cog.Builder<Panel> {
    private Panel internal;

    public PanelBuilder() {
        this.internal = new Panel();
        this.onlyFromThisDashboard(false);
        this.onlyInTimeRange(false);
        this.limit(10);
        this.showUser(true);
        this.showTime(true);
        this.showTags(true);
        this.navigateToPanel(true);
        this.navigateBefore("10m");
        this.navigateAfter("10m");
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
    
    public Panel build() {
        return this.internal;
    }
}
