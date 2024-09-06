package anonymous_struct;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonInclude;

public class SomeStruct {
    @JsonInclude(JsonInclude.Include.NON_NULL)
    @JsonProperty("time")
    public Object time;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private final SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder time(Object time) {
    this.internal.time = time;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
