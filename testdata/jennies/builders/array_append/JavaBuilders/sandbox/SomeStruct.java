package sandbox;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import java.util.LinkedList;

public class SomeStruct { 
    @JsonProperty("tags")
    public List<String> tags;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private final SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder tags(String tags) {
		if (this.internal.tags == null) {
			this.internal.tags = new LinkedList<>();
		}
    this.internal.tags.add(tags);
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
