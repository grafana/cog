package nullable_map_assignment;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonInclude;
import java.util.Map;

public class SomeStruct {
    @JsonInclude(JsonInclude.Include.NON_NULL)
    @JsonProperty("config")
    public Map<String, String> config;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        protected final SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder config(Map<String, String> config) {
    this.internal.config = config;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
