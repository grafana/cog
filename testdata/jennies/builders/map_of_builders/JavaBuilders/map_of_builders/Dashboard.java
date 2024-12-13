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
    public Builder panels(Map<String, cog.Builder<Panel>> panels) {
        Map<String, Panel> panelsResource = new HashMap<>();
        for (var entry : panels.entrySet()) {
           panelsResource.put(entry.getKey(), entry.getValue().build());
        }
    this.internal.panels = panelsResource;
        return this;
    }
    public Dashboard build() {
            return this.internal;
        }
    }
}
