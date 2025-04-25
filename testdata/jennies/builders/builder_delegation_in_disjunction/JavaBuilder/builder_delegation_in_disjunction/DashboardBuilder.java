package builder_delegation_in_disjunction;

import java.util.List;

public class DashboardBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public DashboardBuilder() {
        this.internal = new Dashboard();
    }
    public DashboardBuilder singleLinkOrString(cog.Builder<unknown> singleLinkOrString) {
        this.internal.singleLinkOrString = singleLinkOrString.build();
        return this;
    }
    
    public DashboardBuilder linksOrStrings(cog.Builder<List<unknown>> linksOrStrings) {
        this.internal.linksOrStrings = linksOrStrings.build();
        return this;
    }
    
    public DashboardBuilder disjunctionOfBuilders(cog.Builder<unknown> disjunctionOfBuilders) {
        this.internal.disjunctionOfBuilders = disjunctionOfBuilders.build();
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
