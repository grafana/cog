package sandbox;

import java.util.HashMap;

public class SomeStructWithDefaultEnumBuilder implements cog.Builder<SomeStructWithDefaultEnum> {
    protected final SomeStructWithDefaultEnum internal;
    
    public SomeStructWithDefaultEnumBuilder() {
        this.internal = new SomeStructWithDefaultEnum();
    }
    public SomeStructWithDefaultEnumBuilder data(StringEnumWithDefault key,String value) {
		if (this.internal.data == null) {
			this.internal.data = new HashMap<>();
		}
        this.internal.data.put(key, value);
        return this;
    }
    public SomeStructWithDefaultEnum build() {
        return this.internal;
    }
}
