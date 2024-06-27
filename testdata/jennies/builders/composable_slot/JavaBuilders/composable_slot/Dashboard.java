package composable_slot;

import cog.variants.Dataquery;
import java.util.List;

public class Dashboard {
    public Dataquery target;
    public List<Dataquery> targets;
    
    public static class Builder implements cog.Builder<Dashboard> {
        private Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder Target(cog.Builder<Dataquery> target) {
    this.internal.target = target.Build();
        return this;
    }
    
    public Builder Targets(cog.Builder<List<Dataquery>> targets) {
    this.internal.targets = targets.Build();
        return this;
    }
    public Dashboard Build() {
            return this.internal;
        }
    }
}
