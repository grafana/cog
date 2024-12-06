package sandbox;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonSetter;
import com.fasterxml.jackson.annotation.Nulls;
import java.util.Map;
import java.util.HashMap;

public class SomeStruct {
    @JsonSetter(nulls = Nulls.AS_EMPTY)
    @JsonProperty("annotations")
    public Map<String, String> annotations;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        protected final SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder annotations(String key,String value) {
		if (this.internal.annotations == null) {
			this.internal.annotations = new Hashmap<>();
		}
    this.internal.annotations[key] = value;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
