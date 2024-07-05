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
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder Editable() {
    this.internal.editable = true;
        return this;
    }
    
    public Builder Readonly() {
    this.internal.editable = false;
        return this;
    }
    
    public Builder AutoRefresh() {
    this.internal.autoRefresh = true;
        return this;
    }
    
    public Builder NoAutoRefresh() {
    this.internal.autoRefresh = false;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}
