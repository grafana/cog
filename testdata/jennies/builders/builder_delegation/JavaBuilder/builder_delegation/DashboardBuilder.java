package builder_delegation;

import java.util.List;

public class DashboardBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public DashboardBuilder() {
        this.internal = new Dashboard();
    }
    public DashboardBuilder id(Long id) {
    this.internal.id = id;
        return this;
    }
    
    public DashboardBuilder title(String title) {
    this.internal.title = title;
        return this;
    }
    
    public DashboardBuilder links(cog.Builder<List<DashboardLink>> links) {
    this.internal.links = links.build();
        return this;
    }
    
    public DashboardBuilder linksOfLinks(cog.Builder<List<List<DashboardLink>>> linksOfLinks) {
    this.internal.linksOfLinks = linksOfLinks.build();
        return this;
    }
    
    public DashboardBuilder singleLink(cog.Builder<DashboardLink> singleLink) {
    this.internal.singleLink = singleLink.build();
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
