package sandbox;


public class SomeStruct {
    public Object time;
    
    public static class Builder {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder setTime(String from,String to) {
    	if (this.time == null) {
			this.time = new Object();
		}
    this.internal.time.from = from;
    	if (this.time == null) {
			this.time = new Object();
		}
    this.internal.time.to = to;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
