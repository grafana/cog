package constructor_initializations;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class SomePanel { 
    @JsonProperty("type")
    public String type; 
    @JsonProperty("title")
    public String title; 
    @JsonProperty("cursor")
    public CursorMode cursor;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomePanel> {
        private final SomePanel internal;
        
        public Builder() {
            this.internal = new SomePanel();
    this.internal.type = "panel_type";
    this.internal.cursor = CursorMode.TOOLTIP;
        }
    public Builder title(String title) {
    this.internal.title = title;
        return this;
    }
    public SomePanel build() {
            return this.internal;
        }
    }
}
