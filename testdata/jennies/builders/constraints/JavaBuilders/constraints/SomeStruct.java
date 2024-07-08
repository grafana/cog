package constraints;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class SomeStruct { 
    @JsonProperty("id")
    public Long id; 
    @JsonProperty("title")
    public String title;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private final SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder id(Long id) {
        if (!(id >= 5)) {
            throw new IllegalArgumentException("id must be >= 5");
        }
        if (!(id < 10)) {
            throw new IllegalArgumentException("id must be < 10");
        }
    this.internal.id = id;
        return this;
    }
    
    public Builder title(String title) {
        if (!(title.length() >= 1)) {
            throw new IllegalArgumentException("title.length() must be >= 1");
        }
    this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
