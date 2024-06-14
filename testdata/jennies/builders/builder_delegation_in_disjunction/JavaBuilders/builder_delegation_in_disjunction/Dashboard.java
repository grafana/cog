package builder_delegation_in_disjunction;

import java.util.List;

public class Dashboard {
    // will be expanded to cog.Builder<DashboardLink> | string
    public unknown singleLinkOrString;
    // will be expanded to [](cog.Builder<DashboardLink> | string)
    public List<unknown> linksOrStrings;
    public unknown disjunctionOfBuilders;
        
    public static class Builder {
        private Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder setSingleLinkOrString(cog.Builder<unknown> singleLinkOrString) {
    this.internal.singleLinkOrString = singleLinkOrString.build();
        return this;
    }
    
    public Builder setLinksOrStrings(cog.Builder<List<unknown>> linksOrStrings) {
    this.internal.linksOrStrings = linksOrStrings.build();
        return this;
    }
    
    public Builder setDisjunctionOfBuilders(cog.Builder<unknown> disjunctionOfBuilders) {
    this.internal.disjunctionOfBuilders = disjunctionOfBuilders.build();
        return this;
    }
    public Dashboard build() {
            return this.internal;
        }
    }
}
