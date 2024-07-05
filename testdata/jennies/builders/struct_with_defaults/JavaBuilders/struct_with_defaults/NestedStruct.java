package struct_with_defaults;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class NestedStruct { 
    @JsonProperty("stringVal")
    public String stringVal; 
    @JsonProperty("intVal")
    public Long intVal;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<NestedStruct> {
        private NestedStruct internal;
        
        public Builder() {
            this.internal = new NestedStruct();
        }
    public Builder StringVal(String stringVal) {
    this.internal.stringVal = stringVal;
        return this;
    }
    
    public Builder IntVal(Long intVal) {
    this.internal.intVal = intVal;
        return this;
    }
    public NestedStruct Build() {
            return this.internal;
        }
    }
}
