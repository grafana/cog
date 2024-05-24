package com.grafana.foundation.sandbox;


public class SomeStruct {
    public String title;
    
    public static class Builder {
        private SomeStruct internal;
        
        public Builder(String title) {
            this.internal = new SomeStruct();
    this.internal.title = title;
        }
    public Builder setTitle(String title) {
    this.internal.title = title;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
