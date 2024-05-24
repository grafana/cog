package com.grafana.foundation.struct_with_defaults;


public class NestedStruct {
    public String stringVal;
    public Long intVal;
    
    public static class Builder {
        private NestedStruct internal;
        
        public Builder() {
            this.internal = new NestedStruct();
        }
    public Builder setStringVal(String stringVal) {
    this.internal.stringVal = stringVal;
        return this;
    }
    
    public Builder setIntVal(Long intVal) {
    this.internal.intVal = intVal;
        return this;
    }
    public NestedStruct build() {
            return this.internal;
        }
    }
}
