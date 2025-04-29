package map_of_builders;

import java.util.Map;
import java.util.HashMap;

public class DashboardBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public DashboardBuilder() {
        this.internal = new Dashboard();
    }
    public DashboardBuilder panels(Map<String, cog.Builder<Panel>> panels) {
        Map<String, Panel> panelsResources = new HashMap<>();
        for (var entry1 : panels.entrySet()) {
                panelsDepth1 = entry1.getValue().build();
                panelsResources.put(entry1.getKey(), panelsDepth1);           
        }
        this.internal.panels = panelsResources;
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
