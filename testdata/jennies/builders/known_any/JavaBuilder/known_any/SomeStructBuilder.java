package known_any;


public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder title(String title) {
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
