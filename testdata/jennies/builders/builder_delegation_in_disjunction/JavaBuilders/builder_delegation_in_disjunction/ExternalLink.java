package builder_delegation_in_disjunction;


public class ExternalLink {
    public String url;
    
    public static class Builder implements cog.Builder<ExternalLink> {
        private ExternalLink internal;
        
        public Builder() {
            this.internal = new ExternalLink();
        }
    public Builder Url(String url) {
    this.internal.url = url;
        return this;
    }
    public ExternalLink Build() {
            return this.internal;
        }
    }
}
