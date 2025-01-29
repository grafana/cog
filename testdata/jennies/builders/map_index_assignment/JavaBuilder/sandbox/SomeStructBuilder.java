package sandbox;

import java.util.HashMap;

public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder annotations(String key,String value) {
		if (this.internal.annotations == null) {
			this.internal.annotations = new HashMap<>();
		}
        this.internal.annotations.put(key, value);
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
