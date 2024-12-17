package map_of_builders;

import java.util.Map;

public class DashboardBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public DashboardBuilder() {
        this.internal = new Dashboard();
    }
    public DashboardBuilder panels(cog.Builder<Map<String, Panel>> panels) {
    this.internal.panels = panels.build();
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
