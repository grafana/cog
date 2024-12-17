package map_of_builders;

import java.util.Map;
import java.util.HashMap;

public class DashboardBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public DashboardBuilder() {
        this.internal = new Dashboard();
    }
    public DashboardBuilder panels(Map<String, cog.Builder<Panel>> panels) {
        Map<String, Panel> panelsResource = new HashMap<>();
        for (var entry : panels.entrySet()) {
           panelsResource.put(entry.getKey(), entry.getValue().build());
        }
    this.internal.panels = panelsResource;
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
