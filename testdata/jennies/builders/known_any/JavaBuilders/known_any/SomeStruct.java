package known_any;


public class SomeStruct {
    public Object config;
    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder Title(String title) {
		if (this.internal.config == null) {
			this.internal.config = new known_any.Config();
		}
        known_any.Config configResource = (known_any.Config) this.internal.config;
        configResource.title = title;
    this.internal.config = configResource;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
