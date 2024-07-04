package composable_slot;

import cog.variants.Dataquery;
import java.util.List;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class Dashboard { 
    @JsonProperty("target")
    public Dataquery target; 
    @JsonProperty("targets")
    public List<Dataquery> targets;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
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
