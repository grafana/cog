package dataquery_variant_builder;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class Loki implements cog.variants.Dataquery { 
    @JsonProperty("expr")
    public String expr;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Loki> {
        private final Loki internal;
        
        public Builder() {
            this.internal = new Loki();
        }
    public Builder expr(String expr) {
    this.internal.expr = expr;
        return this;
    }
    public Loki build() {
            return this.internal;
        }
    }
    
}
