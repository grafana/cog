package some_pkg;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;
import com.fasterxml.jackson.annotation.JsonInclude;
import other_pkg.Name;

public class Person {
    @JsonInclude(JsonInclude.Include.NON_EMPTY)
    @JsonProperty("name")
    public Name name;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Person> {
        protected final Person internal;
        
        public Builder() {
            this.internal = new Person();
        }
    public Builder name(Name name) {
    this.internal.name = name;
        return this;
    }
    public Person build() {
            return this.internal;
        }
    }
}
