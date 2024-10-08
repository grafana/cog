package builder_delegation;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonSetter;
import com.fasterxml.jackson.annotation.Nulls;
import java.util.List;
import com.fasterxml.jackson.annotation.JsonInclude;

public class Dashboard {
    @JsonProperty("id")
    public Long id;
    @JsonProperty("title")
    public String title;
    // will be expanded to []cog.Builder<DashboardLink>
    @JsonSetter(nulls = Nulls.AS_EMPTY)
    @JsonProperty("links")
    public List<DashboardLink> links;
    // will be expanded to [][]cog.Builder<DashboardLink>
    @JsonSetter(nulls = Nulls.AS_EMPTY)
    @JsonProperty("linksOfLinks")
    public List<List<DashboardLink>> linksOfLinks;
    // will be expanded to cog.Builder<DashboardLink>
    @JsonInclude(JsonInclude.Include.NON_EMPTY)
    @JsonProperty("singleLink")
    public DashboardLink singleLink;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Dashboard> {
        protected final Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder id(Long id) {
    this.internal.id = id;
        return this;
    }
    
    public Builder title(String title) {
    this.internal.title = title;
        return this;
    }
    
    public Builder links(cog.Builder<List<DashboardLink>> links) {
    this.internal.links = links.build();
        return this;
    }
    
    public Builder linksOfLinks(cog.Builder<List<List<DashboardLink>>> linksOfLinks) {
    this.internal.linksOfLinks = linksOfLinks.build();
        return this;
    }
    
    public Builder singleLink(cog.Builder<DashboardLink> singleLink) {
    this.internal.singleLink = singleLink.build();
        return this;
    }
    public Dashboard build() {
            return this.internal;
        }
    }
}
