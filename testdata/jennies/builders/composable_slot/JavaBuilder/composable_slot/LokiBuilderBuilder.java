package composable_slot;

import cog.variants.Dataquery;
import java.util.List;

public class LokiBuilderBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public LokiBuilderBuilder() {
        this.internal = new Dashboard();
    }
    public LokiBuilderBuilder target(cog.Builder<Dataquery> target) {
        this.internal.target = target.build();
        return this;
    }
    
    public LokiBuilderBuilder targets(cog.Builder<List<Dataquery>> targets) {
        this.internal.targets = targets.build();
        return this;
    }
    public Dashboard build() {
        return this.internal;
    }
}
