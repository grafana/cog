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
        this.internal.target = target.build();
        return this;
    }
    
    public LokiBuilderBuilder targets(List<cog.Builder<Dataquery>> targets) {
        List<Dataquery> targetsResource = new LinkedList<>();
        for (List<Dataquery> targetsVal : targets) {
           targetsResource.add(targetsVal.build());
        }
        this.internal.targets = targetsResource;
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
