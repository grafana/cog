package builder_delegation_in_disjunction;

import java.util.List;
import java.util.LinkedList;

public class DashboardBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public DashboardBuilder() {
        this.internal = new Dashboard();
    }
    public DashboardBuilder singleLinkOrString(cog.Builder<unknown> singleLinkOrString) {
    unknown singleLinkOrStringResource = singleLinkOrString.build();
        this.internal.singleLinkOrString = singleLinkOrStringResource;
        return this;
    }
    
    public DashboardBuilder linksOrStrings(List<cog.Builder<unknown>> linksOrStrings) {
        List<unknown> linksOrStringsResources = new LinkedList<>();
        for (cog.Builder<unknown> r1 : linksOrStrings) {
                unknown linksOrStringsDepth1 = r1.build();
                linksOrStringsResources.add(linksOrStringsDepth1); 
        }
        this.internal.linksOrStrings = linksOrStringsResources;
        return this;
    }
    
    public DashboardBuilder disjunctionOfBuilders(cog.Builder<unknown> disjunctionOfBuilders) {
    unknown disjunctionOfBuildersResource = disjunctionOfBuilders.build();
        this.internal.disjunctionOfBuilders = disjunctionOfBuildersResource;
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
