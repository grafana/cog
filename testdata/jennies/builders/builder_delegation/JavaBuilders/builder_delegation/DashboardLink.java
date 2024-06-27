package builder_delegation;


public class DashboardLink {
    public String title;
    public String url;
    
    public static class Builder implements cog.Builder<DashboardLink> {
        private DashboardLink internal;
        
        public Builder() {
            this.internal = new DashboardLink();
        }
    public Builder Title(String title) {
    this.internal.title = title;
        return this;
    }
    
    public Builder Url(String url) {
    this.internal.url = url;
        return this;
    }
    public DashboardLink Build() {
            return this.internal;
        }
    }
}
