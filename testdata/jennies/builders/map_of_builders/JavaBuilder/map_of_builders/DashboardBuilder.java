package map_of_builders;

import java.util.Map;
import java.util.HashMap;
import java.util.List;

public class DashboardBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public DashboardBuilder() {
        this.internal = new Dashboard();
    }
    public DashboardBuilder panels(Map<String, cog.Builder<Panel>> panels) {
        Map<String, Panel> panelsResources = new HashMap<>();
        for (var entry1 : panels.entrySet()) {
                Panel panelsDepth1 = entry1.getValue().build();
                panelsResources.put(entry1.getKey(), panelsDepth1);           
        }
        this.internal.panels = panelsResources;
        return this;
    }
    
    public DashboardBuilder facets(Map<String, List<cog.Builder<FacetValue>>> facets) {
        Map<String, List<FacetValue>> facetsResources = new HashMap<>();
        for (var entry1 : facets.entrySet()) {
        for (cog.Builder<FacetValue> r2 : entry1) {
                FacetValue facetsDepth2 = r2.build();
                facetsDepth1.add(facetsDepth2); 
        }
                facetsResources.put(entry1.getKey(), entry1.getValue());           
        }
        this.internal.facets = facetsResources;
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
