package builder_delegation_in_disjunction;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class DashboardLink {
    @JsonProperty("title")
    public String title;
    @JsonProperty("url")
    public String url;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<DashboardLink> {
        private final DashboardLink internal;
        
        public Builder() {
            this.internal = new DashboardLink();
        }
    public Builder title(String title) {
    this.internal.title = title;
        return this;
    }
    
    public Builder url(String url) {
    this.internal.url = url;
        return this;
    }
    public DashboardLink build() {
            return this.internal;
        }
    }
}
