package sandbox;

import java.util.List;
import java.util.LinkedList;

public class SomeStruct {
    public List<String> tags;
    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder Tags(String tags) {
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
}
