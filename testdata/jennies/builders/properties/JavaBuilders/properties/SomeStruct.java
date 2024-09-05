package properties;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class SomeStruct {
    @JsonProperty("id")
    public Long id;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private final SomeStruct internal;
        private String someBuilderProperty;
        
        public Builder() {
            this.internal = new SomeStruct();
        this.someBuilderProperty = "";
        }
    public Builder id(Long id) {
    this.internal.id = id;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
