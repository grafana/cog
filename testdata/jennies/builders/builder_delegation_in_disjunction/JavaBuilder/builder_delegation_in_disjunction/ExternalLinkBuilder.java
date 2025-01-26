package builder_delegation_in_disjunction;


public class ExternalLinkBuilder implements cog.Builder<ExternalLink> {
    protected final ExternalLink internal;
    
    public ExternalLinkBuilder() {
        this.internal = new ExternalLink();
    }
    public ExternalLinkBuilder url(String url) {
        this.internal.url = url;
        return this;
    }
    public ExternalLink build() {
        return this.internal;
    }
}
