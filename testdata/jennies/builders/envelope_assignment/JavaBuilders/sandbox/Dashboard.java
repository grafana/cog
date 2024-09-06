package sandbox;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonSetter;
import com.fasterxml.jackson.annotation.Nulls;
import java.util.List;
import java.util.LinkedList;

public class Dashboard {
    @JsonSetter(nulls = Nulls.AS_EMPTY)
    @JsonProperty("variables")
    public List<Variable> variables;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Dashboard> {
        private final Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder withVariable(String name,String value) {
		if (this.internal.variables == null) {
			this.internal.variables = new LinkedList<>();
		}
    Variable variable = new Variable();
        variable.name = name;
        variable.value = value;
    this.internal.variables.add(variable);
        return this;
    }
    public Dashboard build() {
            return this.internal;
        }
    }
}
