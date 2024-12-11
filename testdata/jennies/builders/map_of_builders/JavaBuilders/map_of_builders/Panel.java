package map_of_builders;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class Panel {
    @JsonProperty("title")
    public String title;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Panel> {
        protected final Panel internal;
        
        public Builder() {
            this.internal = new Panel();
        }
    public Builder title(String title) {
    this.internal.title = title;
        return this;
    }
    public Panel build() {
            return this.internal;
        }
    }
}
