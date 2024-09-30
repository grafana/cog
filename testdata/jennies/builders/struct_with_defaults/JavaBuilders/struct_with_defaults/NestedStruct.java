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
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<NestedStruct> {
        protected final NestedStruct internal;
        
        public Builder() {
            this.internal = new NestedStruct();
        }
    public Builder stringVal(String stringVal) {
    this.internal.stringVal = stringVal;
        return this;
    }
    
    public Builder intVal(Long intVal) {
    this.internal.intVal = intVal;
        return this;
    }
    public NestedStruct build() {
            return this.internal;
        }
    }
}
