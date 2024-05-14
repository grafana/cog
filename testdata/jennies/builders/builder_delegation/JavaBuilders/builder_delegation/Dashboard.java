package builder_delegation;

import java.util.List;

public class Dashboard {
    public Long id;
    public String title;
    // will be expanded to []cog.Builder<DashboardLink>
    public List<DashboardLink> links;
    // will be expanded to [][]cog.Builder<DashboardLink>
    public List<List<DashboardLink>> linksOfLinks;
    // will be expanded to cog.Builder<DashboardLink>
    public DashboardLink singleLink;
    
    public static class Builder {
        private Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder setId(Long id) {
    this.internal.id = id;
        return this;
    }
    
    public Builder setTitle(String title) {
    this.internal.title = title;
        return this;
    }
    
    public Builder setLinks(cog.Builder<List<DashboardLink>> links) {
    this.internal.links = links.build();
        return this;
    }
    
    public Builder setLinksOfLinks(cog.Builder<List<List<DashboardLink>>> linksOfLinks) {
    this.internal.linksOfLinks = linksOfLinks.build();
        return this;
    }
    
    public Builder setSingleLink(cog.Builder<DashboardLink> singleLink) {
    this.internal.singleLink = singleLink.build();
        return this;
    }
    public Dashboard build() {
            return this.internal;
        }
    }
}
