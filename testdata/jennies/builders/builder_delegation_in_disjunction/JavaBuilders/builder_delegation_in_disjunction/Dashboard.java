package builder_delegation_in_disjunction;

import java.util.List;

public class Dashboard {
    // will be expanded to cog.Builder<DashboardLink> | string
    public unknown singleLinkOrString;
    // will be expanded to [](cog.Builder<DashboardLink> | string)
    public List<unknown> linksOrStrings;
    public unknown disjunctionOfBuilders;
    
    public static class Builder implements cog.Builder<Dashboard> {
        private Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder SingleLinkOrString(cog.Builder<unknown> singleLinkOrString) {
    this.internal.singleLinkOrString = singleLinkOrString.Build();
        return this;
    }
    
    public Builder LinksOrStrings(cog.Builder<List<unknown>> linksOrStrings) {
    this.internal.linksOrStrings = linksOrStrings.Build();
        return this;
    }
    
    public Builder DisjunctionOfBuilders(cog.Builder<unknown> disjunctionOfBuilders) {
    this.internal.disjunctionOfBuilders = disjunctionOfBuilders.Build();
        return this;
    }
    public Dashboard Build() {
            return this.internal;
        }
    }
}
