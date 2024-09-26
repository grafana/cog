package dashboard;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.ObjectWriter;

public class Panel {
    @JsonProperty("onlyFromThisDashboard")
    public Boolean onlyFromThisDashboard;
    
    public String toJSON() throws JsonProcessingException {
        ObjectWriter ow = new ObjectMapper().writer().withDefaultPrettyPrinter();
        return ow.writeValueAsString(this);
    }

    
    public static class Builder<T extends Builder<T>> implements cog.Builder<Panel> {
        protected final Panel internal;
        
        public Builder() {
            this.internal = new Panel();
        this.onlyFromThisDashboard(false);
        }
    public T onlyFromThisDashboard(Boolean onlyFromThisDashboard) {
    this.internal.onlyFromThisDashboard = onlyFromThisDashboard;
        return (T) this;
    }
    public Panel build() {
            return this.internal;
        }
    }
}
