package sandbox;


public class SomeStructBuilder implements cog.Builder<SomeStruct> {
    protected final SomeStruct internal;
    
    public SomeStructBuilder() {
        this.internal = new SomeStruct();
    }
    public SomeStructBuilder time(String from,String to) {
		if (this.internal.time == null) {
			this.internal.time = new Object();
		}
        this.internal.time.from = from;
        this.internal.time.to = to;
        return this;
    }
    public SomeStruct build() {
        return this.internal;
    }
}
