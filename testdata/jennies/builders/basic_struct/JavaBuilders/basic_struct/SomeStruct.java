package basic_struct;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonSetter;
import com.fasterxml.jackson.annotation.Nulls;
import java.util.List;

// SomeStruct, to hold data.
public class SomeStruct {
    // id identifies something. Weird, right?
    @JsonProperty("id")
    public Long id;
    @JsonProperty("uid")
    public String uid;
    @JsonSetter(nulls = Nulls.AS_EMPTY)
    @JsonProperty("tags")
    public List<String> tags;
    // This thing could be live.
    // Or maybe not.
    @JsonProperty("liveNow")
    public Boolean liveNow;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        protected final SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder id(Long id) {
    this.internal.id = id;
        return this;
    }
    
    public Builder uid(String uid) {
    this.internal.uid = uid;
        return this;
    }
    
    public Builder tags(List<String> tags) {
    this.internal.tags = tags;
        return this;
    }
    
    public Builder liveNow(Boolean liveNow) {
    this.internal.liveNow = liveNow;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
