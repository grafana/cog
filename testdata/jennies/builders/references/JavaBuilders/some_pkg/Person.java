package some_pkg;

import other_pkg.Name;
import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class Person { 
    @JsonProperty("name")
    public Name name;
    
    public String ToJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder implements cog.Builder<Person> {
        private Person internal;
        
        public Builder() {
            this.internal = new Person();
        }
    public Builder Name(Name name) {
    this.internal.name = name;
        return this;
    }
    public Person Build() {
            return this.internal;
        }
    }
}
