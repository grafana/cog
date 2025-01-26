package sandbox;

import java.util.LinkedList;

public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder tags(String tags) {
		if (this.internal.tags == null) {
			this.internal.tags = new LinkedList<>();
		}
        this.internal.tags.add(tags);
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
