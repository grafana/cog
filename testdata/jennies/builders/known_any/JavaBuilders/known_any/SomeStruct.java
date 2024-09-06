package known_any;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonInclude;

public class SomeStruct {
    @JsonInclude(JsonInclude.Include.NON_NULL)
    @JsonProperty("config")
    public Object config;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private final SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder title(String title) {
		if (this.internal.config == null) {
			this.internal.config = new known_any.Config();
		}
        known_any.Config configResource = (known_any.Config) this.internal.config;
        configResource.title = title;
    this.internal.config = configResource;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
