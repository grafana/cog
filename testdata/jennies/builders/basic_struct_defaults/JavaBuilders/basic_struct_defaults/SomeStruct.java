package basic_struct_defaults;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class SomeStruct { 
    @JsonProperty("id")
    public Long id; 
    @JsonProperty("uid")
    public String uid; 
    @JsonProperty("tags")
    public List<String> tags; 
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
        this.Id(42L);
        this.Uid("default-uid");
        this.Tags(List.of("generated", "cog"));
        this.LiveNow(true);
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
