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
        ((Config) this.internal.config).title = title;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
