package nullable_map_assignment;

import java.util.Map;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class SomeStruct { 
    @JsonProperty("config")
    public Map<String, String> config;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder Config(Map<String, String> config) {
    this.internal.config = config;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}
