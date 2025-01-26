package builder_delegation;


public class DashboardLinkBuilder implements cog.Builder<DashboardLink> {
    protected final DashboardLink internal;
    
    public DashboardLinkBuilder() {
        this.internal = new DashboardLink();
    }
    public DashboardLinkBuilder title(String title) {
        this.internal.title = title;
        return this;
    }
    
    public DashboardLinkBuilder url(String url) {
        this.internal.url = url;
        return this;
    }
    public DashboardLink build() {
        return this.internal;
    }
}
