package builder_delegation_in_disjunction;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class ExternalLink { 
    @JsonProperty("url")
    public String url;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<ExternalLink> {
        private ExternalLink internal;
        
        public Builder() {
            this.internal = new ExternalLink();
        }
    public Builder Url(String url) {
    this.internal.url = url;
        return this;
    }
    public ExternalLink Build() {
            return this.internal;
        }
    }
}
