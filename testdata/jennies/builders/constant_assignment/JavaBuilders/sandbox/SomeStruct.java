package sandbox;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class SomeStruct { 
    @JsonProperty("editable")
    public unknown editable; 
    @JsonProperty("autoRefresh")
    public unknown autoRefresh;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private final SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder editable() {
    this.internal.editable = true;
        return this;
    }
    
    public Builder readonly() {
    this.internal.editable = false;
        return this;
    }
    
    public Builder autoRefresh() {
    this.internal.autoRefresh = true;
        return this;
    }
    
    public Builder noAutoRefresh() {
    this.internal.autoRefresh = false;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
