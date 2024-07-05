package sandbox;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class SomeStruct { 
    @JsonProperty("title")
    public String title;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder(String title) {
            this.internal = new SomeStruct();
    this.internal.title = title;
        }
    public Builder Title(String title) {
    this.internal.title = title;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}
