package composable_slot;

import cog.variants.Dataquery;
import java.util.List;
import java.util.LinkedList;

public class LokiBuilderBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public LokiBuilderBuilder() {
        this.internal = new Dashboard();
    }
    public LokiBuilderBuilder target(cog.Builder<Dataquery> target) {
    Dataquery targetResource = target.build();
        this.internal.target = targetResource;
        return this;
    }
    
    public LokiBuilderBuilder targets(List<cog.Builder<Dataquery>> targets) {
        List<Dataquery> targetsResources = new LinkedList<>();
        for (cog.Builder<Dataquery> r1 : targets) {
                Dataquery targetsDepth1 = r1.build();
                targetsResources.add(targetsDepth1); 
        }
        this.internal.targets = targetsResources;
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
