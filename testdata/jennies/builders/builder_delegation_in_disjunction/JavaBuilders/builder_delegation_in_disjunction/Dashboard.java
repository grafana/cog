package builder_delegation_in_disjunction;

import java.util.List;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class Dashboard {
    // will be expanded to cog.Builder<DashboardLink> | string 
    @JsonProperty("singleLinkOrString")
    public unknown singleLinkOrString;
    // will be expanded to [](cog.Builder<DashboardLink> | string) 
    @JsonProperty("linksOrStrings")
    public List<unknown> linksOrStrings; 
    @JsonProperty("disjunctionOfBuilders")
    public unknown disjunctionOfBuilders;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Dashboard> {
        private Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder SingleLinkOrString(cog.Builder<unknown> singleLinkOrString) {
    this.internal.singleLinkOrString = singleLinkOrString.Build();
        return this;
    }
    
    public Builder LinksOrStrings(cog.Builder<List<unknown>> linksOrStrings) {
    this.internal.linksOrStrings = linksOrStrings.Build();
        return this;
    }
    
    public Builder DisjunctionOfBuilders(cog.Builder<unknown> disjunctionOfBuilders) {
    this.internal.disjunctionOfBuilders = disjunctionOfBuilders.Build();
        return this;
    }
    public Dashboard Build() {
            return this.internal;
        }
    }
}
