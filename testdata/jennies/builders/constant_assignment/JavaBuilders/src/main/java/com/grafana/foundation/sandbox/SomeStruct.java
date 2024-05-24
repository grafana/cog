package com.grafana.foundation.sandbox;


public class SomeStruct {
    public unknown editable;
    public unknown autoRefresh;
    
    public static class Builder {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder setEditable() {
    this.internal.editable = true;
        return this;
    }
    
    public Builder setReadonly() {
    this.internal.editable = false;
        return this;
    }
    
    public Builder setAutoRefresh() {
    this.internal.autoRefresh = true;
        return this;
    }
    
    public Builder setNoAutoRefresh() {
    this.internal.autoRefresh = false;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
