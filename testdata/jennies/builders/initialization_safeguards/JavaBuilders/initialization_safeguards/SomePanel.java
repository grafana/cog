package initialization_safeguards;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class SomePanel { 
    @JsonProperty("title")
    public String title; 
    @JsonProperty("options")
    public Options options;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<SomePanel> {
        private SomePanel internal;
        
        public Builder() {
            this.internal = new SomePanel();
        }
    public Builder Title(String title) {
    this.internal.title = title;
        return this;
    }
    
    public Builder ShowLegend(Boolean show) {
		if (this.internal.options == null) {
			this.internal.options = new initialization_safeguards.Options();
		}
		if (this.internal.options.legend == null) {
			this.internal.options.legend = new initialization_safeguards.LegendOptions();
		}
    this.internal.options.legend.show = show;
        return this;
    }
    public SomePanel Build() {
            return this.internal;
        }
    }
}
