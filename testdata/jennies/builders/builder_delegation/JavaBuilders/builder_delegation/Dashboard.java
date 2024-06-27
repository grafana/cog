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
    
    public static class Builder implements cog.Builder<Dashboard> {
        private Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder Id(Long id) {
    this.internal.id = id;
        return this;
    }
    
    public Builder Title(String title) {
    this.internal.title = title;
        return this;
    }
    
    public Builder Links(cog.Builder<List<DashboardLink>> links) {
    this.internal.links = links.Build();
        return this;
    }
    
    public Builder LinksOfLinks(cog.Builder<List<List<DashboardLink>>> linksOfLinks) {
    this.internal.linksOfLinks = linksOfLinks.Build();
        return this;
    }
    
    public Builder SingleLink(cog.Builder<DashboardLink> singleLink) {
    this.internal.singleLink = singleLink.Build();
        return this;
    }
    public Dashboard Build() {
            return this.internal;
        }
    }
}
