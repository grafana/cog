package builder_delegation;

import java.util.List;
import java.util.LinkedList;

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
    
    public DashboardBuilder links(List<cog.Builder<DashboardLink>> links) {
        List<DashboardLink> linksResources = new LinkedList<>();
        for (cog.Builder<DashboardLink> r1 : links) {
                DashboardLink linksDepth1 = r1.build();
                linksResources.add(linksDepth1); 
        }
        this.internal.links = linksResources;
        return this;
    }
    
    public DashboardBuilder linksOfLinks(List<List<cog.Builder<DashboardLink>>> linksOfLinks) {
        List<List<DashboardLink>> linksOfLinksResources = new LinkedList<>();
        for (List<cog.Builder<DashboardLink>> r1 : linksOfLinks) {
                List<DashboardLink> linksOfLinksDepth1 = new LinkedList<>();
        for (cog.Builder<DashboardLink> r2 : r1) {
                DashboardLink linksOfLinksDepth2 = r2.build();
                linksOfLinksDepth1.add(linksOfLinksDepth2); 
        }
                
                linksOfLinksResources.add(linksOfLinksDepth1); 
        }
        this.internal.linksOfLinks = linksOfLinksResources;
        return this;
    }
    
    public DashboardBuilder singleLink(cog.Builder<DashboardLink> singleLink) {
    DashboardLink singleLinkResource = singleLink.build();
        this.internal.singleLink = singleLinkResource;
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
