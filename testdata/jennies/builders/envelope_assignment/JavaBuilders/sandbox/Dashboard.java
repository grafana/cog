package sandbox;

import java.util.List;
import java.util.LinkedList;

public class Dashboard {
    public List<Variable> variables;
    
    public static class Builder {
        private Dashboard internal;
        
        public Builder() {
            this.internal = new Dashboard();
        }
    public Builder setWithVariable(String name,String value) {
		if (this.variables == null) {
			this.variables = new LinkedList<>();
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
