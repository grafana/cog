package basic_struct;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

// SomeStruct, to hold data.
public class SomeStruct {
    // id identifies something. Weird, right? 
    @JsonProperty("id")
    public Long id; 
    @JsonProperty("uid")
    public String uid; 
    @JsonProperty("tags")
    public List<String> tags;
    // This thing could be live.
    // Or maybe not. 
    @JsonProperty("liveNow")
    public Boolean liveNow;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder Id(Long id) {
    this.internal.id = id;
        return this;
    }
    
    public Builder Uid(String uid) {
    this.internal.uid = uid;
        return this;
    }
    
    public Builder Tags(List<String> tags) {
    this.internal.tags = tags;
        return this;
    }
    
    public Builder LiveNow(Boolean liveNow) {
    this.internal.liveNow = liveNow;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}
