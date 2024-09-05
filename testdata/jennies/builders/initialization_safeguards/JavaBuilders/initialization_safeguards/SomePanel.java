package initialization_safeguards;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;

public class SomePanel {
    @JsonProperty("title")
    public String title;
    @JsonSerialize(include = JsonSerialize.Inclusion.NON_NULL)
    @JsonProperty("options")
    public Options options;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomePanel> {
        private final SomePanel internal;
        
        public Builder() {
            this.internal = new SomePanel();
        }
    public Builder title(String title) {
    this.internal.title = title;
        return this;
    }
    
    public Builder showLegend(Boolean show) {
		if (this.internal.options == null) {
			this.internal.options = new initialization_safeguards.Options();
		}
		if (this.internal.options.legend == null) {
			this.internal.options.legend = new initialization_safeguards.LegendOptions();
		}
    this.internal.options.legend.show = show;
        return this;
    }
    public SomePanel build() {
            return this.internal;
        }
    }
}
