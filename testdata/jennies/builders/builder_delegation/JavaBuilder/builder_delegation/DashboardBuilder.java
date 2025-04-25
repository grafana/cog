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
        List<DashboardLink> linksResource = new LinkedList<>();
        for (List<DashboardLink> linksVal : links) {
           linksResource.add(linksVal.build());
        }
        this.internal.links = linksResource;
        return this;
    }
    
    public DashboardBuilder linksOfLinks(List<List<cog.Builder<DashboardLink>>> linksOfLinks) {
        List<List<DashboardLink>> linksOfLinksResource = new LinkedList<>();
        for (List<List<DashboardLink>> linksOfLinksVal : linksOfLinks) {
           linksOfLinksResource.add(linksOfLinksVal.build());
        }
        this.internal.linksOfLinks = linksOfLinksResource;
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
