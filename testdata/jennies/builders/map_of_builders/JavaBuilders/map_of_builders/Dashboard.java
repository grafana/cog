package map_of_builders;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonSetter;
import com.fasterxml.jackson.annotation.Nulls;
import java.util.Map;

public class Dashboard {
    @JsonSetter(nulls = Nulls.AS_EMPTY)
    @JsonProperty("panels")
    public Map<String, Panel> panels;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Dashboard> {
        protected final Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder panels(cog.Builder<Map<String, Panel>> panels) {
    this.internal.panels = panels.build();
        return this;
    }
    public Dashboard build() {
            return this.internal;
        }
    }
}
