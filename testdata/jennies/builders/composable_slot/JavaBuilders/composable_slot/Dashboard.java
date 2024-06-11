package composable_slot;

import cog.variants.Dataquery;
import java.util.List;

public class Dashboard {
    public Dataquery target;
    public List<Dataquery> targets;
    
    public static class Builder {
        private Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder setTarget(cog.Builder<Dataquery> target) {
    this.internal.target = target.build();
        return this;
    }
    
    public Builder setTargets(cog.Builder<List<Dataquery>> targets) {
    this.internal.targets = targets.build();
        return this;
    }
    public Dashboard build() {
            return this.internal;
        }
    }
}
