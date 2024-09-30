package composable_slot;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.databind.annotation.JsonDeserialize;
import cog.variants.Dataquery;
import com.fasterxml.jackson.annotation.JsonSetter;
import com.fasterxml.jackson.annotation.Nulls;
import java.util.List;

@JsonDeserialize(using = DashboardDeserializer.class)
public class Dashboard {
    @JsonProperty("target")
    public Dataquery target;
    @JsonSetter(nulls = Nulls.AS_EMPTY)
    @JsonProperty("targets")
    public List<Dataquery> targets;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Dashboard> {
        protected final Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder target(cog.Builder<Dataquery> target) {
    this.internal.target = target.build();
        return this;
    }
    
    public Builder targets(cog.Builder<List<Dataquery>> targets) {
    this.internal.targets = targets.build();
        return this;
    }
    public Dashboard build() {
            return this.internal;
        }
    }
}
