package sandbox;


public class SomeStruct {
    public Object time;
    
    public static class Builder implements cog.Builder<SomeStruct> {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder Time(String from,String to) {
		if (this.internal.time == null) {
			this.internal.time = new Object();
		}
    this.internal.time.from = from;
    this.internal.time.to = to;
        return this;
    }
    public SomeStruct Build() {
            return this.internal;
        }
    }
}
