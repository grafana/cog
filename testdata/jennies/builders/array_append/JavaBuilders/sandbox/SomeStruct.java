package sandbox;

import java.util.List;
import java.util.LinkedList;

public class SomeStruct {
    public List<String> tags;
    
    public static class Builder {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder setTags(String tags) {
		if (this.tags == null) {
			this.tags = new LinkedList<>();
		}
    this.internal.tags.add(tags);
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
