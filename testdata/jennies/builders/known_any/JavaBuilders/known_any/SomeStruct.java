package known_any;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class SomeStruct { 
    @JsonProperty("config")
    public Object config;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder Title(String title) {
		if (this.internal.config == null) {
			this.internal.config = new known_any.Config();
		}
        known_any.Config configResource = (known_any.Config) this.internal.config;
        configResource.title = title;
    this.internal.config = configResource;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}
