package sandbox;

import java.util.HashMap;

public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder data(StringEnum key,String value) {
		if (this.internal.data == null) {
			this.internal.data = new HashMap<>();
		}
        this.internal.data.put(key, value);
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
