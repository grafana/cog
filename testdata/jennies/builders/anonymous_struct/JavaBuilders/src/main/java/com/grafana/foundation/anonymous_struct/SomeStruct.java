package com.grafana.foundation.anonymous_struct;


public class SomeStruct {
    public Object time;
    
    public static class Builder {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder setTime(Object time) {
    this.internal.time = time;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
