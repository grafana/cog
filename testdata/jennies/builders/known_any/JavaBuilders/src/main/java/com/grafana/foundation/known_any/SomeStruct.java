package com.grafana.foundation.known_any;


public class SomeStruct {
    public Object config;
    
    public static class Builder {
        private SomeStruct internal;
        
        public Builder() {
            this.internal = new SomeStruct();
        }
    public Builder setTitle(String title) {
		if (this.internal.config == null) {
			this.internal.config = new com.grafana.foundation.known_any.Config();
		}
        com.grafana.foundation.known_any.Config configResource = (com.grafana.foundation.known_any.Config) this.internal.config;
        configResource.title = title;
    this.internal.config = configResource;
        return this;
    }
    public SomeStruct build() {
            return this.internal;
        }
    }
}
