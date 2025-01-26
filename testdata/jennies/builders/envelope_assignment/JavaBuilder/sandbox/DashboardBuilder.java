package sandbox;

import java.util.LinkedList;

public class DashboardBuilder implements cog.Builder<Dashboard> {
    protected final Dashboard internal;
    
    public DashboardBuilder() {
        this.internal = new Dashboard();
    }
    public DashboardBuilder withVariable(String name,String value) {
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
